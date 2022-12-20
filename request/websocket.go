package request

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/braycarlson/asol/authorization"
)

type (
	Websocket struct {
		authorization *authorization.Authorization
	}
)

func (websocket *Websocket) Credential() string {
	username := websocket.authorization.Username
	password := websocket.authorization.Password

	return base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
}

func (websocket *Websocket) WebsocketAddress() string {
	port := websocket.authorization.Port
	return "wss://127.0.0.1:" + port
}

func (websocket *Websocket) LocalAddress() string {
	port := websocket.authorization.Port
	return "https://127.0.0.1:" + port
}

func (websocket *Websocket) Get(uri string) (*http.Request, error) {
	uri = websocket.LocalAddress() + uri
	request, err := http.NewRequest(http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	return request, nil
}

func (websocket *Websocket) Post(uri string, data []byte) (*http.Request, error) {
	uri = websocket.LocalAddress() + uri
	buffer := bytes.NewBuffer(data)

	request, err := http.NewRequest(http.MethodPost, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	return request, nil
}

func (websocket *Websocket) Patch(uri string, data []byte) (*http.Request, error) {
	uri = websocket.LocalAddress() + uri
	buffer := bytes.NewBuffer(data)

	request, err := http.NewRequest(http.MethodPatch, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	return request, nil
}

func (websocket *Websocket) Put(uri string, data []byte) (*http.Request, error) {
	uri = websocket.LocalAddress() + uri
	buffer := bytes.NewBuffer(data)

	request, err := http.NewRequest(http.MethodPut, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	return request, nil
}

func (websocket *Websocket) Delete(uri string) (*http.Request, error) {
	uri = websocket.LocalAddress() + uri
	request, err := http.NewRequest(http.MethodDelete, uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	return request, nil
}

func (websocket *Websocket) Request(client *HTTPClient, request *http.Request) ([]byte, error) {
	client.setInsecureSkipVerify(true)

	request.Header.Set(
		"Authorization",
		"Basic "+websocket.Credential(),
	)

	response, err := client.client.Do(request)

	if err != nil {
		return nil, &ClientError{"WebsocketRequest", err}
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	return bytes, err
}
