package wshub

// broker listener
// broker sender
// broker

// broker listener

// listens for message from broker
type BrokerListener[Message any] interface {
	listen(listener func(Message))
}

type brokerListenerImpl[Message any] struct {
	listeners []func(m Message)
}

// nolint:U1000 // Ignore unused function warning.
func (impl *brokerListenerImpl[Message]) listen(listener func(Message)) {
	impl.listeners = append(impl.listeners, listener)
}

func NewBrokerListener[Message any]() (brokerListener BrokerListener[Message], invoke func(m Message)) {
	listeners := &brokerListenerImpl[Message]{
		listeners: []func(m Message){},
	}

	return listeners, func(m Message) {
		for _, listener := range listeners.listeners {
			listener(m)
		}
	}
}

// broker sender

// sends message to broker
type BrokerSender[Message any] interface {
	send(message Message)
}

type brokerSenderImpl[Message any] struct {
	receive func(message Message)
}

// nolint:U1000 // Ignore unused function warning.
func (impl *brokerSenderImpl[Message]) send(message Message) {
	impl.receive(message)
}

func NewBrokerSender[Message any](receive func(m Message)) BrokerSender[Message] {
	return &brokerSenderImpl[Message]{receive: receive}
}

// broker

type broker struct {
	// is sent by broker when started to ensure all sockets are removed
	started BrokerSender[Started]
	// sends to broker after receiving
	connectRequest BrokerSender[ConnectRequest]
	// connects web socket
	connectConfirmation BrokerListener[ConnectConfirmation]

	// sends to broker after receiving socket message
	received BrokerSender[SocketMessage]

	// after being received sends this to user
	respond BrokerListener[SocketMessage]

	// after being received closes socket
	close BrokerListener[Close]
	// sends to broker after closing socket
	closed BrokerSender[Close]
}

func NewBroker(
	started BrokerSender[Started],
	connectRequest BrokerSender[ConnectRequest],
	connectConfirmation BrokerListener[ConnectConfirmation],
	received BrokerSender[SocketMessage],
	respond BrokerListener[SocketMessage],
	close BrokerListener[Close],
	closed BrokerSender[Close],
) *broker {
	return &broker{
		started,
		connectRequest,
		connectConfirmation,
		received,
		respond,
		close,
		closed,
	}
}
