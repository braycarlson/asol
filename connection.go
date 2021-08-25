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

type (
	Asol struct {
		Connection *websocket.Conn
		*GameProcess
		*HTTPClient
		*ConnectionEventManager
		*WebsocketEventManager
		mutex     *sync.Mutex
		reconnect bool
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

func NewAsol() *Asol {
	return &Asol{
		&websocket.Conn{},
		NewGameProcess(),
		NewHTTPClient(),
		&ConnectionEventManager{},
		&WebsocketEventManager{},
		&sync.Mutex{},
		true,
	}
}

func (asol *Asol) refresh() {
	asol.GameProcess = NewGameProcess()
}

func (asol *Asol) respawn(path string) error {
	err := exec.Command(path).Start()
	return err
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

		if login.Connected && strings.ToLower(login.State) == "succeeded" &&
			login.AccountId != 0 && login.SummonerId != 0 {
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

		asol.onOpen()
		asol.isReady()
		asol.onReady()
		asol.isLoggedIn()
		asol.onLogin()

		path := asol.Path()

		asol.listen()

		if asol.reconnect == false {
			break
		}

		for {
			var closed, opened bool

			outer := make(chan bool, 1)

			IsLoginOrClient(
				outer,
				"RiotClientUx.exe",
				"RiotClientServices.exe",
			)

			select {
			case logout := <-outer:
				switch logout {
				case true:
					asol.onLogout()

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
							opened = true
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
				asol.onClientClose()
				asol.respawn(path)
			}

			if opened {
				break
			}

			close(outer)
			time.Sleep(1000 * time.Millisecond)
		}

		asol.refresh()
		asol.onReconnect()
	}
}

func (asol *Asol) listen() {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	connection, _, err := dialer.Dial(
		asol.WebsocketAddress(),
		http.Header{
			"Content-Type":  []string{"application/json"},
			"Accept":        []string{"application/json"},
			"Authorization": {"Basic " + asol.Authorization()},
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
			asol.onWebsocketError(
				fmt.Errorf("%v", err),
			)

			if err == io.ErrUnexpectedEOF {
				continue
			}

			if websocket.IsUnexpectedCloseError(err) {
				asol.onWebsocketClose()
			}

			break
		}

		err = asol.Match(
			&Message{
				URI:    response.data["uri"].(string),
				Method: response.data["eventType"].(string),
				Data:   response.data,
			},
		)

		if err != nil {
			asol.onWebsocketError(
				fmt.Errorf("%v", err),
			)
		}
	}
}
