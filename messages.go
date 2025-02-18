package wshub

// started
// connect confirmation
// connect request
// message
// close

// started

type Started struct{}

func NewStarted() Started {
	return Started{}
}

// connect confirmation

type ConnectConfirmation struct {
	SocketId   Id   `json:"socket_id"`
	CanConnect bool `json:"can_connect"`
}

func NewConnectConfirmation(socketId Id, canConnect bool) ConnectConfirmation {
	return ConnectConfirmation{
		SocketId:   socketId,
		CanConnect: canConnect,
	}
}

// connect request

type ConnectRequest struct {
	SocketId Id     `json:"socket_id"`
	Payload  []byte `json:"info"`
}

func NewConnectRequest(id Id, payload []byte) ConnectRequest {
	return ConnectRequest{SocketId: id, Payload: payload}
}

// message

type SocketMessage struct {
	SocketId Id     `json:"socket_id"`
	Payload  []byte `json:"payload"`
}

func NewSocketMessage(id Id, payload []byte) SocketMessage {
	return SocketMessage{
		SocketId: id,
		Payload:  payload,
	}
}

// close

type Close struct {
	SocketId Id `json:"socket_id"`
}

func NewClose(id Id) Close {
	return Close{
		SocketId: id,
	}
}
