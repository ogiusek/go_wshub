package wshub

import (
	"encoding/json"

	"github.com/google/uuid"
)

// id
// socket conn

// id

type id struct{ value string }

func (u id) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.value)
}

func (u *id) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &u.value)
}

func NewSocketId() id {
	return id{value: uuid.NewString()}
}

func SocketIdFrom(value string) id {
	return id{value: value}
}

// socket conn

type SocketConn interface {
	Close()
	ReadMessage() ([]byte, error)
	Send([]byte)
}
