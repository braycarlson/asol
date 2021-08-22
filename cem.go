package asol

type (
	EventCallback  func()
	WebsocketError func(error)

	ConnectionEventManager struct {
		onOpen           EventCallback
		onReady          EventCallback
		onLogin          EventCallback
		onLogout         EventCallback
		onClientClose    EventCallback
		onWebsocketClose EventCallback
		onReconnect      EventCallback
		onWebsocketError WebsocketError
	}
)

func (asol *Asol) OnOpen(callback EventCallback) {
	asol.ConnectionEventManager.onOpen = callback
}

func (asol *Asol) OnReady(callback EventCallback) {
	asol.ConnectionEventManager.onReady = callback
}

func (asol *Asol) OnLogin(callback EventCallback) {
	asol.ConnectionEventManager.onLogin = callback
}

func (asol *Asol) OnLogout(callback EventCallback) {
	asol.ConnectionEventManager.onLogout = callback
}

func (asol *Asol) OnClientClose(callback EventCallback) {
	asol.ConnectionEventManager.onClientClose = callback
}

func (asol *Asol) OnWebsocketClose(callback EventCallback) {
	asol.ConnectionEventManager.onWebsocketClose = callback
}

func (asol *Asol) OnReconnect(callback EventCallback) {
	asol.ConnectionEventManager.onReconnect = callback
}

func (asol *Asol) OnWebsocketError(callback WebsocketError) {
	asol.ConnectionEventManager.onWebsocketError = callback
}
