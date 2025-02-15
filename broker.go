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
	listeners []func(Message)
}

func (impl *brokerListenerImpl[Message]) listen(listener func(Message)) {
	impl.listeners = append(impl.listeners, listener)
}

func NewBrokerListener[Message any]() (brokerListener BrokerListener[Message], invoke func(Message)) {
	listeners := &brokerListenerImpl[Message]{
		listeners: []func(Message){},
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

func (impl *brokerSenderImpl[Message]) send(message Message) {
	impl.receive(message)
}

func NewBrokerSender[Message any](receive func(Message)) BrokerSender[Message] {
	return &brokerSenderImpl[Message]{receive: receive}
}

// broker

type broker struct {
	// is sent by broker when started to ensure all sockets are removed
	Started BrokerSender[Started]
	// sends to broker after receiving
	ConnectRequest BrokerSender[ConnectRequest]
	// connects web socket
	ConnectConfirmation BrokerListener[ConnectConfirmation]

	// sends to broker after receiving socket message
	Received BrokerSender[SocketMessage]

	// after being received sends this to user
	Respond BrokerListener[SocketMessage]

	// after being received closes socket
	Close BrokerListener[Close]
	// sends to broker after closing socket
	Closed BrokerSender[Close]
}

func NewBroker(
	Started BrokerSender[Started],
	ConnectRequest BrokerSender[ConnectRequest],
	ConnectConfirmation BrokerListener[ConnectConfirmation],
	Received BrokerSender[SocketMessage],
	Respond BrokerListener[SocketMessage],
	Close BrokerListener[Close],
	Closed BrokerSender[Close],
) *broker {
	return &broker{
		Started,
		ConnectRequest,
		ConnectConfirmation,
		Received,
		Respond,
		Close,
		Closed,
	}
}
