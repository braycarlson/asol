package main

import (
	"log"

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

func onOpen() {
	log.Println("The client is opened")
}

func onReady() {
	log.Println("The client is ready")
}

func onLogin() {
	log.Println("The client is logged in")
}

func onLogout() {
	log.Println("The client is logged out")
}

func onClientClose() {
	log.Println("The client is closed")
}

func onWebsocketClose() {
	log.Println("The client's websocket closed")
}

func onReconnect() {
	log.Println("The client is reconnected")
}

func onWebsocketError(error error) {
	log.Println(error)
}

func onCollection(message []byte) {
	log.Println(string(message))
}

func onGame(message []byte) {
	log.Println(message)
}

func main() {
	client.OnOpen(onOpen)
	client.OnReady(onReady)
	client.OnLogin(onLogin)
	client.OnLogout(onLogout)
	client.OnClientClose(onClientClose)
	client.OnWebsocketClose(onWebsocketClose)
	client.OnReconnect(onReconnect)
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
