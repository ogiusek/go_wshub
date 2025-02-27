package wshub

import (
	"encoding/json"

	"github.com/google/uuid"
)

// id
// socket conn

// id

type Id struct{ value string }

func (u Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.value)
}

func (u *Id) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &u.value)
}

func (u *Id) String() string {
	return u.value
}

func NewSocketId() Id {
	return Id{value: uuid.NewString()}
}

func SocketIdFrom(value string) Id {
	return Id{value: value}
}

// socket conn

type SocketConn interface {
	Close()
	ReadMessage() ([]byte, error)
	Send([]byte)
}
