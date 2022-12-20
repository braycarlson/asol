package cem

type (
	EventCallback  func()
	ProcessError   func(error)
	SearchError    func(error)
	WebsocketError func(error)

	ConnectionEventManager struct {
		OnSearchCallback         EventCallback
		OnOpenCallback           EventCallback
		OnReadyCallback          EventCallback
		OnLoginCallback          EventCallback
		OnProcessErrorCallback   ProcessError
		OnSearchErrorCallback    SearchError
		OnWebsocketCloseCallback EventCallback
		OnWebsocketErrorCallback WebsocketError
	}
)

func (cem *ConnectionEventManager) OnSearch(callback EventCallback) {
	cem.OnSearchCallback = callback
}

func (cem *ConnectionEventManager) OnOpen(callback EventCallback) {
	cem.OnOpenCallback = callback
}

func (cem *ConnectionEventManager) OnReady(callback EventCallback) {
	cem.OnReadyCallback = callback
}

func (cem *ConnectionEventManager) OnLogin(callback EventCallback) {
	cem.OnLoginCallback = callback
}

func (cem *ConnectionEventManager) OnProcessError(callback ProcessError) {
	cem.OnProcessErrorCallback = callback
}

func (cem *ConnectionEventManager) OnSearchError(callback SearchError) {
	cem.OnSearchErrorCallback = callback
}

func (cem *ConnectionEventManager) OnWebsocketClose(callback EventCallback) {
	cem.OnWebsocketCloseCallback = callback
}

func (cem *ConnectionEventManager) OnWebsocketError(callback WebsocketError) {
	cem.OnWebsocketErrorCallback = callback
}
