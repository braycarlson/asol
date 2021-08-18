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
	"time"
)

type (
	RiotClient struct {
		*http.Client
	}

	WebClient struct {
		*http.Client
	}
)

func NewRiotClient() *RiotClient {
	return &RiotClient{
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

func NewWebClient() *WebClient {
	return &WebClient{
		&http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 5 * time.Second,
					DualStack: true,
				}).DialContext,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
				DisableKeepAlives:     true,
			},
		},
	}
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

func (asol *Asol) WebsocketHeader() http.Header {
	authorization := asol.Authorization()

	return http.Header{
		"Content-Type":  []string{"application/json"},
		"Accept":        []string{"application/json"},
		"Connection":    []string{"close"},
		"Authorization": {"Basic " + authorization},
	}
}

func (asol *Asol) WebHeader() http.Header {
	return http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
		"Connection":   []string{"close"},
	}
}

func (asol *Asol) NewWebRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	request.Header = asol.WebHeader()
	return request, nil
}

func (asol *Asol) NewGetRequest(endpoint string) (*http.Request, error) {
	uri := asol.LocalAddress() + endpoint
	request, err := http.NewRequest(http.MethodGet, uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header = asol.WebsocketHeader()
	return request, nil
}

func (asol *Asol) NewPostRequest(endpoint string, data map[string]interface{}) (*http.Request, error) {
	uri := asol.LocalAddress() + endpoint
	payload, err := json.Marshal(data)
	buffer := bytes.NewBuffer(payload)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header = asol.WebsocketHeader()
	return request, nil
}

func (asol *Asol) NewPatchRequest(endpoint string, data map[string]interface{}) (*http.Request, error) {
	uri := asol.LocalAddress() + endpoint
	payload, err := json.Marshal(data)
	buffer := bytes.NewBuffer(payload)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPatch, uri, buffer)

	if err != nil {
		return nil, err
	}

	request.Header = asol.WebsocketHeader()
	return request, nil
}

func (asol *Asol) NewDeleteRequest(endpoint string) (*http.Request, error) {
	uri := asol.LocalAddress() + endpoint
	request, err := http.NewRequest(http.MethodDelete, uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header = asol.WebsocketHeader()
	return request, nil
}

func (asol *Asol) RawRequest(request *http.Request) ([]byte, error) {
	response, err := asol.RiotClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(
		io.LimitReader(response.Body, 1048576),
	)

	return bytes, err
}

func (asol *Asol) RiotRequest(request *http.Request) (interface{}, error) {
	var data interface{}
	response, err := asol.RiotClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&data)
	return data, err
}

func (asol *Asol) WebRequest(request *http.Request) (interface{}, error) {
	var data interface{}
	response, err := asol.WebClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&data)
	return data, err
}
