package asol

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type (
	Client struct {
		*http.Client
	}
)

func NewClient() *Client {
	return &Client{
		&http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 5 * time.Second,
					DualStack: true,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
				DisableKeepAlives:     true,
			},
		},
	}
}

func (client *Client) setInsecureSkipVerify(insecureSkipVerify bool) {
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

func (asol *Asol) Header() http.Header {
	return http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
}

func (asol *Asol) Get(uri string) (*http.Request, error) {
	if strings.HasPrefix(uri, "/") {
		uri = asol.LocalAddress() + uri
	}

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	request.Header = asol.Header()

	if err != nil {
		return nil, err
	}

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
	request.Header = asol.Header()

	if err != nil {
		return nil, err
	}

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
	request.Header = asol.Header()

	if err != nil {
		return nil, err
	}

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
	request.Header = asol.Header()

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (asol *Asol) Delete(uri string) (*http.Request, error) {
	if strings.HasPrefix(uri, "/") {
		uri = asol.LocalAddress() + uri
	}

	request, err := http.NewRequest(http.MethodDelete, uri, nil)
	request.Header = asol.Header()

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (asol *Asol) WebRequest(request *http.Request) (interface{}, error) {
	asol.Client.setInsecureSkipVerify(false)
	response, err := asol.Client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var data interface{}

	json.NewDecoder(response.Body).Decode(&data)
	return data, err
}

func (asol *Asol) RiotRequest(request *http.Request) ([]byte, error) {
	asol.Client.setInsecureSkipVerify(true)

	authorization := asol.Authorization()
	request.Header.Set("Authorization", "Basic "+authorization)

	response, err := asol.Client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(
		io.LimitReader(response.Body, 1048576),
	)

	return bytes, err
}

func (asol *Asol) RawWebRequest(request *http.Request) ([]byte, error) {
	asol.Client.setInsecureSkipVerify(false)
	response, err := asol.Client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(
		io.LimitReader(response.Body, 1048576),
	)

	return bytes, err
}
