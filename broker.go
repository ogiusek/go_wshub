package wshub

// broker listener
// broker sender
// broker

// broker listener

// listens for message from broker
type BrokerListener[Message any] interface {
	Listen(listener func(Message))
}

type brokerListenerImpl[Message any] struct {
	listen func(listener func(Message))
}

func (impl *brokerListenerImpl[Message]) Listen(listener func(Message)) {
	impl.listen(listener)
}

func NewBrokerListener[Message any](listen func(listener func(Message))) BrokerListener[Message] {
	return &brokerListenerImpl[Message]{listen}
}

// broker sender

// sends message to broker
type BrokerSender[Message any] interface {
	Send(message Message)
}

type brokerSenderImpl[Message any] struct {
	send func(message Message)
}

func (impl *brokerSenderImpl[Message]) Send(message Message) {
	impl.send(message)
}

func NewBrokerSender[Message any](send func(Message)) BrokerSender[Message] {
	return &brokerSenderImpl[Message]{send}
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
