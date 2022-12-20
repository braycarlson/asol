package asol

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/braycarlson/asol/authorization"
	"github.com/braycarlson/asol/cem"
	"github.com/braycarlson/asol/game"
	"github.com/braycarlson/asol/request"
	"github.com/braycarlson/asol/wem"
	"github.com/gorilla/websocket"
)

type (
	Asol struct {
		*cem.ConnectionEventManager
		*wem.WebsocketEventManager

		client     *request.HTTPClient
		connection *websocket.Conn
		search     *game.Search
		game       *game.Game
		mutex      *sync.Mutex
		status     bool
	}

	Login struct {
		AccountId      float64
		Connected      bool
		Error          bool
		GasToken       string
		IdToken        string
		IsInLoginQueue bool
		IsNewPlayer    bool
		Puuid          string
		State          string
		SummonerId     float64
		UserAuthToken  string
		Username       string
	}
)

func (login *Login) isReady() bool {
	var state string = strings.ToLower(login.State)

	if login.Connected == true &&
		state == "succeeded" &&
		login.AccountId != 0 &&
		login.SummonerId != 0 {
		return true
	}

	return false
}

func NewAsol() *Asol {
	return &Asol{
		&cem.ConnectionEventManager{},
		&wem.WebsocketEventManager{},

		request.NewHTTPClient(),
		&websocket.Conn{},
		&game.Search{},
		nil,
		&sync.Mutex{},
		false,
	}
}

func (asol *Asol) Client() *request.HTTPClient {
	return asol.client
}

func (asol *Asol) Game() *game.Game {
	return asol.game
}

func (asol *Asol) isGameRunning() bool {
	if asol.game == nil {
		return false
	}

	return true
}

func (asol *Asol) isRunning() bool {
	return asol.status
}

func (asol *Asol) setSearch(search *game.Search) {
	asol.search = search
}

func (asol *Asol) setStatus(status bool) {
	asol.status = status
}

func (asol *Asol) setGame(game *game.Game) {
	asol.game = game
}

func (asol *Asol) isReady() {
	for {
		asol.client.SetWebsocket()
		request, _ := asol.client.Get("/riotclient/region-locale")
		_, err := asol.client.Request(request)

		if err == nil {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (asol *Asol) isLoggedIn() {
	for {
		asol.client.SetWebsocket()
		request, _ := asol.client.Get("/lol-login/v1/session")
		data, err := asol.client.Request(request)

		if err != nil {
			continue
		}

		var login Login
		json.Unmarshal(data, &login)

		if login.isReady() {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (asol *Asol) Start() {
	asol.OnSearchCallback()

	var search *game.Search = game.NewSearch()
	asol.setSearch(search)

	process, err := asol.search.Start()

	if err == nil {
		var game *game.Game = game.NewGame(process)
		var authorization *authorization.Authorization = game.Authorization()

		asol.client.SetAuthorization(authorization)
		asol.setGame(game)

		asol.setStatus(true)

		if asol.game == nil {
			asol.setStatus(false)
			asol.OnProcessErrorCallback(nil)
			return
		}

		err := asol.Registered()

		if err != nil {
			asol.OnWebsocketErrorCallback(err)
			return
		}

		asol.OnOpenCallback()
		asol.isReady()
		asol.OnReadyCallback()
		asol.isLoggedIn()
		asol.OnLoginCallback()

		asol.listen()
	}

	if _, ok := err.(*game.SearchCancelled); ok {
		asol.OnSearchErrorCallback(err)
	}

	if _, ok := err.(*game.ProcessNotFoundError); ok {
		asol.OnProcessErrorCallback(err)
	}

	asol.setStatus(false)
	return
}

func (asol *Asol) Stop() {
	asol.search.Cancel()

	if !asol.isGameRunning() {
		return
	}

	asol.setStatus(false)

	message := []interface{}{wem.Unsubscribe, "OnJsonApiEvent"}
	asol.connection.WriteJSON(&message)

	_, _, err := asol.connection.ReadMessage()

	if err != nil {
		fmt.Errorf("%v", err)
	}

	asol.connection.WriteControl(
		websocket.CloseMessage,
		[]byte{},
		time.Now().Add(time.Second),
	)

	for {
		_, _, err := asol.connection.ReadMessage()

		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return
		}

		if err != nil {
			fmt.Errorf("%v", err)
			return
		}
	}

	asol.setGame(nil)
}

func (asol *Asol) listen() {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	connection, _, err := dialer.Dial(
		asol.client.WebsocketAddress(),
		http.Header{
			"Content-Type":  []string{"application/json"},
			"Accept":        []string{"application/json"},
			"Authorization": {"Basic " + asol.client.Credential()},
		},
	)

	if err != nil {
		asol.OnWebsocketErrorCallback(
			fmt.Errorf("%v", err),
		)

		return
	}

	asol.connection = connection

	asol.mutex.Lock()
	defer asol.mutex.Unlock()

	message := []interface{}{wem.Subscribe, "OnJsonApiEvent"}
	asol.connection.WriteJSON(&message)

	_, _, err = asol.connection.ReadMessage()

	if err != nil {
		asol.OnWebsocketErrorCallback(
			fmt.Errorf("%v", err),
		)
	}

	asol.read()
}

func (asol *Asol) read() {
	defer asol.connection.Close()

	for {
		if asol.isRunning() == false {
			asol.OnWebsocketCloseCallback()
			break
		}

		var response wem.Response
		err := asol.connection.ReadJSON(&response)

		if err != nil {
			if err == io.ErrUnexpectedEOF {
				continue
			}

			asol.OnWebsocketErrorCallback(
				fmt.Errorf("%v", err),
			)

			asol.setStatus(false)
			break
		}

		err = asol.Match(
			&wem.Message{
				URI:    response.Data["uri"].(string),
				Method: response.Data["eventType"].(string),
				Data:   response.Data,
			},
		)

		if err != nil {
			asol.OnWebsocketErrorCallback(
				fmt.Errorf("%v", err),
			)
		}
	}
}
