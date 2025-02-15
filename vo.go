package wshub

import "github.com/google/uuid"

// id
// socket conn

// id

type id struct{ value string }

func NewSocketId() id {
	return id{value: uuid.NewString()}
}

// socket conn

type SocketConn interface {
	Close()
	ReadMessage() ([]byte, error)
	Send([]byte)
}
