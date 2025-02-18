package wshub

import "sync"

// awaiter

type awaiter[Awaited any] struct {
	mux      *sync.Mutex
	awaiters map[Id]chan Awaited
}

func newAwaiter[Awaited any]() *awaiter[Awaited] {
	return &awaiter[Awaited]{
		mux:      &sync.Mutex{},
		awaiters: map[Id]chan Awaited{},
	}
}

func (awaiter *awaiter[Awaited]) Await(id Id) (Awaited, error) {
	awaiter.mux.Lock()
	if _, ok := awaiter.awaiters[id]; ok {
		var nothing Awaited
		return nothing, ErrIdIsTaken
	}
	awaiter.awaiters[id] = make(chan Awaited, 1)
	awaiter.mux.Unlock()
	res := <-awaiter.awaiters[id]
	return res, nil
}

func (awaiter *awaiter[Awaited]) Resolve(id Id, res Awaited) error {
	awaiter.mux.Lock()
	resolver, ok := awaiter.awaiters[id]
	if !ok {
		return ErrIdIsMissing
	}
	delete(awaiter.awaiters, id)
	awaiter.mux.Unlock()

	resolver <- res

	return nil
}
