package wshub

import (
	"sync"
)

// socket storage

type socketStorage struct {
	mux              *sync.Mutex
	sockets          map[Id]SocketConn
	connectListeners []func(Id)
	messageListeners []func(Id, []byte)
	closeListeners   []func(Id)
}

func newSocketStorage() *socketStorage {
	return &socketStorage{
		mux:              &sync.Mutex{},
		sockets:          map[Id]SocketConn{},
		connectListeners: []func(Id){},
		messageListeners: []func(Id, []byte){},
		closeListeners:   []func(Id){},
	}
}

func (storage *socketStorage) run(id Id, conn SocketConn) error {
	defer conn.Close()

	storage.mux.Lock()
	if _, ok := storage.sockets[id]; ok {
		return ErrIdIsTaken
	}
	storage.sockets[id] = conn
	storage.mux.Unlock()

	for _, listener := range storage.connectListeners {
		listener(id)
	}

	for {
		bytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		for _, listener := range storage.messageListeners {
			listener(id, bytes)
		}
	}

	storage.mux.Lock()
	delete(storage.sockets, id)
	storage.mux.Unlock()

	for _, listener := range storage.closeListeners {
		listener(id)
	}

	return nil
}

func (sockets *socketStorage) OnConnect(listener func(id Id)) {
	sockets.connectListeners = append(sockets.connectListeners, listener)
}

func (sockets *socketStorage) OnMessage(listener func(id Id, payload []byte)) {
	sockets.messageListeners = append(sockets.messageListeners, listener)
}

func (sockets *socketStorage) OnClose(listener func(id Id)) {
	sockets.closeListeners = append(sockets.closeListeners, listener)
}

func (storage *socketStorage) Get(id Id) (conn SocketConn, ok bool) {
	_conn, _ok := storage.sockets[id]
	return _conn, _ok
}
