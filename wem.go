package asol

import (
	"encoding/json"
	"fmt"
)

const (
	Welcome     MessageType = 0
	Prefix      MessageType = 1
	Call        MessageType = 2
	CallResult  MessageType = 3
	CallError   MessageType = 4
	Subscribe   MessageType = 5
	Unsubscribe MessageType = 6
	Publish     MessageType = 7
	Event       MessageType = 8
)

type (
	MessageType float64

	WebsocketCallback func([]byte)

	Message struct {
		URI    string
		Method string
		Data   map[string]interface{}
	}

	Response struct {
		messageType MessageType
		event       string
		data        map[string]interface{}
	}

	WebsocketEventManager struct {
		registered []map[string]interface{}
	}
)

func (wem *WebsocketEventManager) Registered() error {
	if len(wem.registered) == 0 {
		return &NoRegisteredEventError{}
	}

	return nil
}

func (wem *WebsocketEventManager) setRegistered(event map[string]interface{}) {
	wem.registered = append(wem.registered, event)
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

func (asol *Asol) Match(message *Message) {
	for _, listener := range asol.registered {
		if message.URI == listener["uri"] && message.Method == listener["method"] {
			callback := listener["callback"].(WebsocketCallback)
			response, err := json.Marshal(message.Data)

			if err != nil {
				asol.onWebsocketError(
					fmt.Errorf("%v", err),
				)
			}

			callback(response)
		}
	}
}

func (response *Response) UnmarshalJSON(message []byte) error {
	return json.Unmarshal(
		message,
		&[]interface{}{
			&response.messageType,
			&response.event,
			&response.data,
		})
}
