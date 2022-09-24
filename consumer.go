package whisper

// Consumer allows to consume messages from a topic.
type Consumer interface {
	// Consume reads the latest message from the topic.
	// If the topic is closed or the broker is stopped
	// Consume returns false
	Consume() (Message, bool)

	// Closed returns true if the consumer poll
	// channel is closed. This is required for consumer groups
	// where multiple topics write onto the same channel
	Closed() bool
}

type consumer struct {
	poll   chan Message
	closed bool
}

func (c *consumer) Consume() (Message, bool) {
	msg, ok := <-c.poll
	if !ok {
		return Message{}, false
	}
	return msg, ok
}

func (c *consumer) close() {
	c.closed = true
}

func (c *consumer) Closed() bool {
	return c.closed
}
