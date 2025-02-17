package wshub

import (
	"log"
)

type WsHub interface {
	// Connect runs until connection is closed/rejected
	Connect(approve func() SocketConn, payload []byte)
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

	broker.started.send(Started{})
	broker.connectConfirmation.listen(func(connectConfirmation ConnectConfirmation) {
		err := awaiter.Resolve(connectConfirmation.SocketId, connectConfirmation)
		if err == ErrIdIsMissing {
			broker.closed.send(NewClose(connectConfirmation.SocketId))
		}
	})

	broker.respond.listen(func(m SocketMessage) {
		conn, ok := storage.Get(m.SocketId)
		if !ok {
			return
		}
		conn.Send(m.Payload)
	})

	broker.close.listen(func(c Close) {
		socket, ok := storage.Get(c.SocketId)
		if !ok {
			return
		}
		socket.Close()
	})

	storage.OnClose(func(id id) {
		broker.closed.send(NewClose(id))
	})

	storage.OnMessage(func(id id, payload []byte) {
		broker.received.send(NewSocketMessage(id, payload))
	})

	return &wshub{
		broker:        broker,
		awaiter:       awaiter,
		socketStorage: storage,
	}
}

func (app *wshub) Connect(resolveConn func() SocketConn, connectRequestPayload []byte) {
	id := NewSocketId()

	app.broker.connectRequest.send(ConnectRequest{
		SocketId: id,
		Payload:  connectRequestPayload,
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
