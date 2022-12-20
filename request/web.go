package request

import (
	"bytes"
	"io"
	"net/http"
)

type (
	Web struct{}
)

func (web *Web) Get(uri string) (*http.Request, error) {
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

func (web *Web) Post(uri string, data []byte) (*http.Request, error) {
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

func (web *Web) Patch(uri string, data []byte) (*http.Request, error) {
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

func (web *Web) Put(uri string, data []byte) (*http.Request, error) {
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

func (web *Web) Delete(uri string) (*http.Request, error) {
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

func (web *Web) Request(client *HTTPClient, request *http.Request) ([]byte, error) {
	client.setInsecureSkipVerify(false)
	response, err := client.client.Do(request)

	if err != nil {
		return nil, &ClientError{"HTTPRequest", err}
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	return bytes, err
}
