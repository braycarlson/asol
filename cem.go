package asol

type (
	EventCallback   func(*Asol)
	RequestCallback func(string, int)
	WebsocketError  func(error)

	ConnectionEventManager struct {
		onOpen           func(*Asol)
		onReady          func(*Asol)
		onLogin          func(*Asol)
		onLogout         func(*Asol)
		onClientClose    func(*Asol)
		onWebsocketClose func(*Asol)
		onReconnect      func(*Asol)
		onRequest        func(string, int)
		onRequestError   func(string, int)
		onWebsocketError func(error)
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

func (asol *Asol) OnRequest(callback RequestCallback) {
	asol.ConnectionEventManager.onRequest = callback
}

func (asol *Asol) OnRequestError(callback RequestCallback) {
	asol.ConnectionEventManager.onRequestError = callback
}

func (asol *Asol) OnWebsocketError(callback WebsocketError) {
	asol.ConnectionEventManager.onWebsocketError = callback
}
