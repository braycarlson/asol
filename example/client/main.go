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

func (client *Client) onSearch() {
	log.Println("The client is searching")
}

func (client *Client) onOpen() {
	log.Println("The client is opened")
}

func (client *Client) onReady() {
	log.Println("The client is ready")
}

func (client *Client) onLogin() {
	log.Println("The client is logged in")
}

func (client *Client) onWebsocketClose() {
	log.Println("The client's websocket closed")
}

func (client *Client) onProcessError(error error) {
	log.Println(error)
}

func (client *Client) onSearchError(error error) {
	log.Println(error)
}

func (client *Client) onWebsocketError(error error) {
	log.Println(error)
}

func (client *Client) onCollection(message []byte) {
	log.Println(string(message))
}

func (client *Client) onGame(message []byte) {
	log.Println(string(message))
}

func main() {
	client := &Client{
		asol.NewAsol(),
	}

	client.OnSearch(client.onSearch)
	client.OnOpen(client.onOpen)
	client.OnReady(client.onReady)
	client.OnLogin(client.onLogin)
	client.OnProcessError(client.onProcessError)
	client.OnSearchError(client.onSearchError)
	client.OnWebsocketClose(client.onWebsocketClose)
	client.OnWebsocketError(client.onWebsocketError)

	client.OnMessage(
		"/lol-settings/v1/account/lol-collection-champions",
		"Update",
		client.onCollection,
	)

	client.OnMessage(
		"/lol-champ-select/v1/asol",
		"Update",
		client.onGame,
	)

	client.Start()
}
