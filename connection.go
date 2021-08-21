package asol

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Asol struct {
	Connection *websocket.Conn
	*GameProcess
	*Client
	*ConnectionEventManager
	*WebsocketEventManager
	mutex     *sync.Mutex
	reconnect bool
}

func NewAsol() *Asol {
	return &Asol{
		&websocket.Conn{},
		NewGameProcess(),
		NewClient(),
		&ConnectionEventManager{},
		&WebsocketEventManager{},
		&sync.Mutex{},
		true,
	}
}

func (asol *Asol) refresh() {
	asol.GameProcess = NewGameProcess()
}

func (asol *Asol) respawn(path string) {
	err := exec.Command(path).Start()

	if err != nil {
		return
	}
}

func (asol *Asol) isReady() {
	for {
		request, _ := asol.Get("/riotclient/region-locale")
		_, err := asol.RiotRequest(request)

		if err == nil {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (asol *Asol) isLoggedIn() {
	for {
		request, _ := asol.Get("/lol-login/v1/session")
		data, err := asol.RiotRequest(request)

		if err != nil {
			continue
		}

		var login Login
		json.Unmarshal(data, &login)

		connected := login.Connected
		state := strings.ToLower(login.State)
		accountId := login.AccountId
		summonerId := login.SummonerId

		if connected && state == "succeeded" && accountId != 0 && summonerId != 0 {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (asol *Asol) Start() {
	for {
		err := asol.Registered()

		if err != nil {
			asol.onWebsocketError(err)
			break
		}

		asol.onOpen(asol)
		asol.isReady()
		asol.onReady(asol)
		asol.isLoggedIn()
		asol.onLogin(asol)

		var path string = asol.Path()
		asol.listen()

		if asol.reconnect == false {
			break
		}

		for {
			closed := false
			game := false

			outer := make(chan bool, 1)

			IsGameOrClient(
				outer,
				"RiotClientUx.exe",
				"RiotClientServices.exe",
			)

			select {
			case logout := <-outer:
				switch logout {
				case true:
					asol.onLogout(asol)

					inner := make(chan bool, 1)

					IsGameOrClient(
						inner,
						"LeagueClientUx.exe",
						"RiotClientServices.exe",
					)

					select {
					case login := <-inner:
						switch login {
						case true:
							game = true
						case false:
							closed = true
						}
					}

					close(inner)
				case false:
					closed = true
				}
			}

			if closed {
				asol.onClientClose(asol)
				asol.respawn(path)
			}

			if game {
				break
			}

			close(outer)
			time.Sleep(1000 * time.Millisecond)
		}

		asol.refresh()
		asol.onReconnect(asol)
	}
}

func (asol *Asol) listen() {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	address := asol.WebsocketAddress()
	authorization := asol.Authorization()

	connection, _, err := dialer.Dial(
		address,
		http.Header{
			"Content-Type":  []string{"application/json"},
			"Accept":        []string{"application/json"},
			"Authorization": {"Basic " + authorization},
		},
	)

	if err != nil {
		asol.onWebsocketError(
			fmt.Errorf("%v", err),
		)

		return
	}

	asol.Connection = connection

	asol.mutex.Lock()

	message := []interface{}{Subscribe, "OnJsonApiEvent"}
	asol.Connection.WriteJSON(&message)

	_, _, err = asol.Connection.ReadMessage()

	if err != nil {
		asol.onWebsocketError(
			fmt.Errorf("%v", err),
		)
	}

	asol.mutex.Unlock()

	asol.read()
}

func (asol *Asol) read() {
	defer asol.Connection.Close()

	for {
		var response Response
		err := asol.Connection.ReadJSON(&response)

		if err != nil {
			if err == io.ErrUnexpectedEOF {
				continue
			}

			asol.onWebsocketError(
				fmt.Errorf("%v", err),
			)

			if websocket.IsUnexpectedCloseError(err) {
				asol.onWebsocketClose(asol)
			}

			break
		}

		asol.Match(
			&Message{
				URI:    response.data["uri"].(string),
				Method: response.data["eventType"].(string),
				Data:   response.data,
			},
		)
	}
}
