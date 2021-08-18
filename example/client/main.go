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

func onOpen(asol *asol.Asol) {
	fmt.Println("Opened")
}

func onReady(asol *asol.Asol) {
	fmt.Println("Ready")
}

func onLogin(asol *asol.Asol) {
	fmt.Println("Logged in")
}

func onLogout(asol *asol.Asol) {
	fmt.Println("Logged out")
}

func onClientClose(asol *asol.Asol) {
	fmt.Println("Client closed")
}

func onWebsocketClose(asol *asol.Asol) {
	fmt.Println("Websocket closed")
}

func onReconnect(asol *asol.Asol) {
	fmt.Println("Reconnected")
}

func onError(error error) {
	fmt.Println(error)
}

func onCollection(asol *asol.Asol, message *asol.Message) {
	fmt.Println(message.Data)
}

func onGame(asol *asol.Asol, message *asol.Message) {
	fmt.Println(message.Data)
}

func main() {
	client := NewClient()

	client.OnOpen(onOpen)
	client.OnReady(onReady)
	client.OnLogin(onLogin)
	client.OnLogout(onLogout)
	client.OnClientClose(onClientClose)
	client.OnWebsocketClose(onWebsocketClose)
	client.OnReconnect(onReconnect)
	client.OnError(onError)

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
