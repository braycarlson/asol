package request

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/braycarlson/asol/authorization"
)

type (
	HTTPClient struct {
		client        *http.Client
		authorization *authorization.Authorization
		strategy      RequestStrategy
	}

	ClientError struct {
		message string
		error   error
	}

	RequestStrategy interface {
		Request(*HTTPClient, *http.Request) ([]byte, error)
		Get(string) (*http.Request, error)
		Post(string, []byte) (*http.Request, error)
		Patch(string, []byte) (*http.Request, error)
		Put(string, []byte) (*http.Request, error)
		Delete(string) (*http.Request, error)
	}
)

func NewHTTPClient() *HTTPClient {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.DialContext = (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
		DualStack: true,
	}).DialContext

	transport.DisableKeepAlives = false
	transport.ExpectContinueTimeout = 5 * time.Second
	transport.MaxIdleConns = 25
	transport.MaxConnsPerHost = 25
	transport.MaxIdleConnsPerHost = 25
	transport.ResponseHeaderTimeout = 5 * time.Second
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	transport.TLSHandshakeTimeout = 5 * time.Second

	return &HTTPClient{
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: transport,
		},
		authorization: &authorization.Authorization{},
		strategy:      &Websocket{},
	}
}

func (error *ClientError) Error() string {
	return fmt.Sprintf("%s: %v", error.message, error.error)
}

func (client *HTTPClient) SetWeb() {
	client.strategy = &Web{}
}

func (client *HTTPClient) SetWebsocket() {
	client.strategy = &Websocket{
		authorization: client.authorization,
	}
}

func (client *HTTPClient) setInsecureSkipVerify(insecureSkipVerify bool) {
	tlsClientConfig := &tls.Config{InsecureSkipVerify: insecureSkipVerify}
	client.client.Transport.(*http.Transport).TLSClientConfig = tlsClientConfig
}

func (client *HTTPClient) SetAuthorization(authorization *authorization.Authorization) {
	client.authorization = authorization
}

func (client *HTTPClient) Credential() string {
	username := client.authorization.Username
	password := client.authorization.Password

	return base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
}

func (client *HTTPClient) WebsocketAddress() string {
	port := client.authorization.Port
	return "wss://127.0.0.1:" + port
}

func (client *HTTPClient) LocalAddress() string {
	port := client.authorization.Port
	return "https://127.0.0.1:" + port
}

func (client *HTTPClient) Get(uri string) (*http.Request, error) {
	return client.strategy.Get(uri)
}

func (client *HTTPClient) Post(uri string, data []byte) (*http.Request, error) {
	return client.strategy.Post(uri, data)
}

func (client *HTTPClient) Patch(uri string, data []byte) (*http.Request, error) {
	return client.strategy.Patch(uri, data)
}

func (client *HTTPClient) Put(uri string, data []byte) (*http.Request, error) {
	return client.strategy.Put(uri, data)
}

func (client *HTTPClient) Delete(uri string) (*http.Request, error) {
	return client.strategy.Delete(uri)
}

func (client *HTTPClient) Request(request *http.Request) ([]byte, error) {
	return client.strategy.Request(client, request)
}
