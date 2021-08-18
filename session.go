package asol

import (
	"crypto/tls"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Asol struct {
	Connection *websocket.Conn
	*GameProcess
	*RiotClient
	*WebClient
	*ConnectionEventManager
	*WebsocketEventManager
	mutex     sync.Mutex
	reconnect bool
}

func NewAsol() *Asol {
	return &Asol{
		&websocket.Conn{},
		NewGameProcess(),
		NewRiotClient(),
		NewWebClient(),
		&ConnectionEventManager{},
		&WebsocketEventManager{},
		sync.Mutex{},
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

func (asol *Asol) poll(endpoint string) {
	for {
		request, _ := asol.NewGetRequest(endpoint)
		_, err := asol.RiotRequest(request)

		if err == nil {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (asol *Asol) Start() {
	for {
		err := asol.Registered()

		if err != nil {
			asol.onError(err)
			break
		}

		asol.onOpen(asol)
		asol.poll("/riotclient/region-locale")
		asol.onReady(asol)
		asol.poll("/lol-login/v1/session")
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
	header := asol.WebsocketHeader()

	connection, _, err := dialer.Dial(address, header)

	if err != nil {
		asol.onError(
			fmt.Errorf("%v", err),
		)

		return
	}

	asol.Connection = connection

	message := []interface{}{Subscribe, "OnJsonApiEvent"}

	asol.mutex.Lock()
	asol.Connection.WriteJSON(&message)
	asol.mutex.Unlock()

	_, _, err = asol.Connection.ReadMessage()

	if err != nil {
		asol.onError(
			fmt.Errorf("%v", err),
		)
	}

	asol.read()

	asol.Connection.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"",
		),
	)
}

func (asol *Asol) read() {
	defer asol.Connection.Close()

	for {
		var response Response
		err := asol.Connection.ReadJSON(&response)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				asol.onWebsocketClose(asol)

				break
			}

			asol.onError(
				fmt.Errorf("%v", err),
			)

			continue
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
