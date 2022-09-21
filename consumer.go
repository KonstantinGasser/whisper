package ares

type Consumer interface {
	Consume() (interface{}, bool)
}

type consumer struct {
	stream <-chan interface{}
}

func (c consumer) Consume() (interface{}, bool) {
	data, ok := <-c.stream
	return data, ok
}
