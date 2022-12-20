package wem

import (
	"encoding/json"
	"fmt"
)

const (
	Welcome     float64 = 0
	Prefix      float64 = 1
	Call        float64 = 2
	CallResult  float64 = 3
	CallError   float64 = 4
	Subscribe   float64 = 5
	Unsubscribe float64 = 6
	Publish     float64 = 7
	Event       float64 = 8
)

type (
	WebsocketCallback func([]byte)

	Message struct {
		URI    string
		Method string
		Data   map[string]interface{}
	}

	Response struct {
		MessageType float64
		Event       string
		Data        map[string]interface{}
	}

	WebsocketEventManager struct {
		registered []map[string]interface{}
	}

	NoRegisteredEventError struct{}
)

func (error *NoRegisteredEventError) Error() string {
	return fmt.Sprintf("No event(s) registered.")
}

func (wem *WebsocketEventManager) Registered() error {
	if len(wem.registered) == 0 {
		return &NoRegisteredEventError{}
	}

	return nil
}

func (wem *WebsocketEventManager) setRegistered(event map[string]interface{}) {
	wem.registered = append(wem.registered, event)
}

func (wem *WebsocketEventManager) Match(message *Message) error {
	for _, listener := range wem.registered {
		if message.URI == listener["uri"] && message.Method == listener["method"] {
			callback := listener["callback"].(WebsocketCallback)
			response, err := json.Marshal(message.Data)

			if err != nil {
				return err
			}

			callback(response)
		}
	}

	return nil
}

func (wem *WebsocketEventManager) OnMessage(uri string, method string, callback WebsocketCallback) WebsocketCallback {
	event := map[string]interface{}{
		"uri":      uri,
		"method":   method,
		"callback": callback,
	}

	wem.setRegistered(event)
	return callback
}

func (response *Response) UnmarshalJSON(message []byte) error {
	return json.Unmarshal(
		message,
		&[]interface{}{
			&response.MessageType,
			&response.Event,
			&response.Data,
		})
}
