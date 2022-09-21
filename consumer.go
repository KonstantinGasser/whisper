package whisper

// Consumer allows to consume messages from a topic.
type Consumer interface {
	// Consume reads the latest message from the topic.
	// If the topic is closed or the broker is stopped
	// Consume returns false
	Consume() (Message, bool)
}

type consumer struct {
	poll <-chan Message
}

func (c consumer) Consume() (Message, bool) {
	msg, ok := <-c.poll
	return msg, ok
}
