package asol

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type (
	HTTPClient struct {
		*http.Client
	}

	ClientError struct {
		message string
		error   error
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
		&http.Client{
			Timeout:   time.Second * 10,
			Transport: transport,
		},
	}
}

func (error *ClientError) Error() string {
	return fmt.Sprintf("%s: %v", error.message, error.error)
}

func (client *HTTPClient) setInsecureSkipVerify(insecureSkipVerify bool) {
	tlsClientConfig := &tls.Config{InsecureSkipVerify: insecureSkipVerify}
	client.Client.Transport.(*http.Transport).TLSClientConfig = tlsClientConfig
}

func (asol *Asol) Authorization() string {
	username := asol.Username()
	password := asol.Password()

	return base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
}

func (asol *Asol) WebsocketAddress() string {
	port := asol.Port()
	return "wss://127.0.0.1:" + port
}

func (asol *Asol) LocalAddress() string {
	port := asol.Port()
	return "https://127.0.0.1:" + port
}

func (asol *Asol) Get(uri string) (*http.Request, error) {
	if strings.HasPrefix(uri, "/") {
		uri = asol.LocalAddress() + uri
	}

	request, err := http.NewRequest(http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request, nil
}

func (asol *Asol) Post(uri string, data map[string]interface{}) (*http.Request, error) {
	if strings.HasPrefix(uri, "/") {
		uri = asol.LocalAddress() + uri
	}

	payload, err := json.Marshal(data)
	buffer := bytes.NewBuffer(payload)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request, nil
}

func (asol *Asol) Patch(uri string, data map[string]interface{}) (*http.Request, error) {
	if strings.HasPrefix(uri, "/") {
		uri = asol.LocalAddress() + uri
	}

	payload, err := json.Marshal(data)
	buffer := bytes.NewBuffer(payload)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPatch, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request, nil
}

func (asol *Asol) Put(uri string, data map[string]interface{}) (*http.Request, error) {
	if strings.HasPrefix(uri, "/") {
		uri = asol.LocalAddress() + uri
	}

	payload, err := json.Marshal(data)
	buffer := bytes.NewBuffer(payload)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPut, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request, nil
}

func (asol *Asol) Delete(uri string) (*http.Request, error) {
	if strings.HasPrefix(uri, "/") {
		uri = asol.LocalAddress() + uri
	}

	request, err := http.NewRequest(http.MethodDelete, uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request, nil
}

func (asol *Asol) RiotRequest(request *http.Request) ([]byte, error) {
	asol.HTTPClient.setInsecureSkipVerify(true)

	request.Header.Set(
		"Authorization",
		"Basic "+asol.Authorization(),
	)

	response, err := asol.HTTPClient.Do(request)

	if err != nil {
		return nil, &ClientError{"RiotRequest", err}
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(
		io.LimitReader(response.Body, 1048576),
	)

	return bytes, err
}

func (asol *Asol) WebRequest(request *http.Request) ([]byte, error) {
	asol.HTTPClient.setInsecureSkipVerify(false)
	response, err := asol.HTTPClient.Do(request)

	if err != nil {
		return nil, &ClientError{"WebRequest", err}
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(
		io.LimitReader(response.Body, 1048576),
	)

	return bytes, err
}
