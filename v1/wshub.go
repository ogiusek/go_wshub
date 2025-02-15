package wshub

import (
	"encoding/json"
	"log"
)

type WsHub interface {
	Connect(resolveConn func() SocketConn, payload any)
}

type wshub struct {
	broker        *broker
	awaiter       *awaiter[ConnectConfirmation]
	socketStorage *socketStorage
}

// connect

func NewWsHub(broker *broker) WsHub {
	storage := newSocketStorage()
	awaiter := newAwaiter[ConnectConfirmation]()

	broker.ConnectConfirmation.Listen(func(connectConfirmation ConnectConfirmation) {
		err := awaiter.Resolve(connectConfirmation.SocketId, connectConfirmation)
		if err == ErrIdIsMissing {
			broker.Closed.Send(NewClose(connectConfirmation.SocketId))
		}
	})

	broker.Respond.Listen(func(m SocketMessage) {
		conn, ok := storage.Get(m.SocketId)
		if !ok {
			return
		}
		payload, _ := json.Marshal(m.Payload)
		conn.Send(payload)
	})

	broker.Close.Listen(func(c Close) {
		socket, ok := storage.Get(c.SocketId)
		if !ok {
			return
		}
		socket.Close()
	})

	storage.OnClose(func(id id) {
		broker.Closed.Send(NewClose(id))
	})

	storage.OnMessage(func(id id, payload []byte) {
		broker.Received.Send(NewSocketMessage(id, payload))
	})

	return &wshub{
		broker:        broker,
		awaiter:       awaiter,
		socketStorage: storage,
	}
}

func (app *wshub) Connect(resolveConn func() SocketConn, payload any) {
	id := NewSocketId()

	app.broker.ConnectRequest.Send(ConnectRequest{
		SocketId: id,
		Payload:  payload,
	})

	connectionConfirmation, err := app.awaiter.Await(id)
	if err != nil {
		log.Panic("uuid is not random")
	}

	if !connectionConfirmation.CanConnect {
		return
	}

	conn := resolveConn()
	app.socketStorage.run(id, conn)
}
