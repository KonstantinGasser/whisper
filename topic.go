package whisper

import (
	"sync"
)

type Message struct {
	Timestamp int64
	Data      interface{}
}

// Topic allows to direct messages to specific consumer groups
// which can subscribe to the topic
type Topic struct {
	sync.RWMutex
	consumer []chan Message
	close    <-chan struct{}
}

func newTopic(route string, opts ...func(*Topic)) *Topic {
	return &Topic{
		consumer: make([]chan Message, 0),
	}
}

func (t *Topic) publish(msg Message) {
	t.RLock()
	defer t.RUnlock()

	for _, c := range t.consumer {
		c <- msg
	}

}

func (t *Topic) subscribe() Consumer {

	pollChan := make(chan Message)

	c := consumer{
		poll: pollChan,
	}

	t.RLock()
	defer t.RUnlock()
	t.consumer = append(t.consumer, pollChan)

	return &c
}

func (t *Topic) stop() {
	t.Lock()
	defer t.Unlock()

	for _, c := range t.consumer {
		close(c)
	}
}
