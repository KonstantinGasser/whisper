package whisper

type Consumer interface {
	Consume() (Message, bool)
}

type consumer struct {
	poll <-chan Message
}

func (c consumer) Consume() (Message, bool) {
	msg, ok := <-c.poll
	return msg, ok
}
