package wshub

// started
// connect confirmation
// connect request
// message
// close

// started

type Started struct{}

// connect confirmation

type ConnectConfirmation struct {
	SocketId   id   `json:"socket_id"`
	CanConnect bool `json:"can_connect"`
}

// connect request

type ConnectRequest struct {
	SocketId id  `json:"socket_id"`
	Payload  any `json:"info"`
}

func NewConnectRequest(id id, payload any) ConnectRequest {
	return ConnectRequest{SocketId: id, Payload: payload}
}

// message

type SocketMessage struct {
	SocketId id  `json:"socket_id"`
	Payload  any `json:"payload"`
}

func NewSocketMessage(id id, payload []byte) SocketMessage {
	return SocketMessage{
		SocketId: id,
		Payload:  payload,
	}
}

// close

type Close struct {
	SocketId id `json:"socket_id"`
}

func NewClose(id id) Close {
	return Close{
		SocketId: id,
	}
}
