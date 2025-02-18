package wshub

import (
	"sync"
)

// socket storage

type socketStorage struct {
	mux              *sync.Mutex
	sockets          map[id]SocketConn
	connectListeners []func(id)
	messageListeners []func(id, []byte)
	closeListeners   []func(id)
}

func newSocketStorage() *socketStorage {
	return &socketStorage{
		mux:              &sync.Mutex{},
		sockets:          map[id]SocketConn{},
		connectListeners: []func(id){},
		messageListeners: []func(id, []byte){},
		closeListeners:   []func(id){},
	}
}

func (storage *socketStorage) run(id id, conn SocketConn) error {
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

func (sockets *socketStorage) OnConnect(listener func(id id)) {
	sockets.connectListeners = append(sockets.connectListeners, listener)
}

func (sockets *socketStorage) OnMessage(listener func(id id, payload []byte)) {
	sockets.messageListeners = append(sockets.messageListeners, listener)
}

func (sockets *socketStorage) OnClose(listener func(id id)) {
	sockets.closeListeners = append(sockets.closeListeners, listener)
}

func (storage *socketStorage) Get(id id) (conn SocketConn, ok bool) {
	_conn, _ok := storage.sockets[id]
	return _conn, _ok
}
