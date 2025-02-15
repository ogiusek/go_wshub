package wshub

// broker listener
// broker sender
// broker

// broker listener

// listens for message from broker
type brokerListener[Message any] interface {
	Listen(listener func(Message))
}

type brokerListenerImpl[Message any] struct {
	listen func(listener func(Message))
}

func (impl *brokerListenerImpl[Message]) Listen(listener func(Message)) {
	impl.listen(listener)
}

func NewBrokerListener[Message any](listen func(listener func(Message))) brokerListener[Message] {
	return &brokerListenerImpl[Message]{listen}
}

// broker sender

// sends message to broker
type brokerSender[Message any] interface {
	Send(message Message)
}

type brokerSenderImpl[Message any] struct {
	send func(message Message)
}

func (impl *brokerSenderImpl[Message]) Send(message Message) {
	impl.send(message)
}

func NewBrokerSender[Message any](send func(Message)) brokerSender[Message] {
	return &brokerSenderImpl[Message]{send}
}

// broker

type broker struct {
	// is sent by broker when started to ensure all sockets are removed
	Started brokerSender[Started]
	// sends to broker after receiving
	ConnectRequest brokerSender[ConnectRequest]
	// connects web socket
	ConnectConfirmation brokerListener[ConnectConfirmation]

	// sends to broker after receiving socket message
	Received brokerSender[SocketMessage]

	// after being received sends this to user
	Respond brokerListener[SocketMessage]

	// after being received closes socket
	Close brokerListener[Close]
	// sends to broker after closing socket
	Closed brokerSender[Close]
}

func NewBroker(
	Started brokerSender[Started],
	ConnectRequest brokerSender[ConnectRequest],
	ConnectConfirmation brokerListener[ConnectConfirmation],
	Received brokerSender[SocketMessage],
	Respond brokerListener[SocketMessage],
	Close brokerListener[Close],
	Closed brokerSender[Close],
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
