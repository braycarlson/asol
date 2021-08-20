package main

import (
	"fmt"

	"github.com/braycarlson/asol"
)

type Client struct {
	*asol.Asol
}

func NewClient() *Client {
	return &Client{
		asol.NewAsol(),
	}
}

var (
	client = &Client{
		asol.NewAsol(),
	}
)

func onOpen(asol *asol.Asol) {
	fmt.Println("The client is opened")
}

func onReady(asol *asol.Asol) {
	fmt.Println("The client is ready")
}

func onLogin(asol *asol.Asol) {
	fmt.Println("The client is logged in")
}

func onLogout(asol *asol.Asol) {
	fmt.Println("The client is logged out")
}

func onClientClose(asol *asol.Asol) {
	fmt.Println("The client is closed")
}

func onWebsocketClose(asol *asol.Asol) {
	fmt.Println("The client's websocket closed")
}

func onReconnect(asol *asol.Asol) {
	fmt.Println("The client is reconnected")
}

func onRequest(uri string, status int) {
	fmt.Println(
		fmt.Sprintf("%d: %s", status, uri),
	)
}

func onRequestError(uri string, status int) {
	fmt.Println(
		fmt.Sprintf("%d: %s", status, uri),
	)
}

func onWebsocketError(error error) {
	fmt.Println(error)
}

func onCollection(asol *asol.Asol, message *asol.Message) {
	fmt.Println(message.Data)
}

func onGame(asol *asol.Asol, message *asol.Message) {
	fmt.Println(message.Data)
}

func main() {
	client.OnOpen(onOpen)
	client.OnReady(onReady)
	client.OnLogin(onLogin)
	client.OnLogout(onLogout)
	client.OnClientClose(onClientClose)
	client.OnWebsocketClose(onWebsocketClose)
	client.OnReconnect(onReconnect)
	client.OnRequest(onRequest)
	client.OnRequestError(onRequestError)
	client.OnWebsocketError(onWebsocketError)

	client.OnMessage(
		"/lol-settings/v1/account/lol-collection-champions",
		"Update",
		onCollection,
	)

	client.OnMessage(
		"/lol-champ-select/v1/asol",
		"Update",
		onGame,
	)

	client.Start()
}
