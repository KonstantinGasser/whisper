package whisper

import (
	"sync"
)

type Message struct {
	Topic     string
	Timestamp int64
	Data      interface{}
}

// Topic allows to direct messages to specific consumer groups
// which can subscribe to the topic
type Topic struct {
	route string
	// mutex used to protect the consumer
	// slice while adding/removing
	sync.RWMutex
	consumer []chan<- Message
	close    <-chan struct{}
}

func newTopic(route string, opts ...func(*Topic)) *Topic {
	return &Topic{
		route:    route,
		consumer: make([]chan<- Message, 0),
	}
}

func (t *Topic) publish(msg Message) {
	msg.Topic = t.route

	t.RLock()
	defer t.RUnlock()

	for _, c := range t.consumer {
		c <- msg
	}

}

// subscribe creates the actual Consumer from which messages
// to the topic can be consumed
func (t *Topic) subscribe(msgChan chan<- Message) {
	t.RLock()
	defer t.RUnlock()

	t.consumer = append(t.consumer, msgChan)
}

func (t *Topic) stop() {
	t.Lock()
	defer t.Unlock()

	for _, c := range t.consumer {
		close(c)
	}
}
