package whisper

import (
	"fmt"
	"sync"
)

type Message struct {
	topic     string
	Timestamp int64
	Data      interface{}
}

func (msg Message) Topic() string { return msg.topic }

// Topic allows to direct messages to specific consumer groups
// which can subscribe to the topic
type Topic struct {
	route string
	// mutex used to protect the consumer
	// slice while adding/removing
	sync.RWMutex
	consumer []*consumer
	close    <-chan struct{}
}

func newTopic(route string, opts ...func(*Topic)) *Topic {
	return &Topic{
		route:    route,
		consumer: make([]*consumer, 0),
	}
}

func (t *Topic) publish(msg Message) {
	msg.topic = t.route

	t.RLock()
	defer t.RUnlock()

	for _, c := range t.consumer {
		c.poll <- msg
	}

}

// subscribe creates the actual Consumer from which messages
// to the topic can be consumed
func (t *Topic) subscribe(cns *consumer) {
	t.RLock()
	defer t.RUnlock()

	t.consumer = append(t.consumer, cns)
}

func (t *Topic) stop() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic in topic: %s: %v\n", t.route, err)
		}
	}()

	t.Lock()
	defer t.Unlock()

	for _, c := range t.consumer {
		if !c.Closed() {
			c.close()
			close(c.poll)
		}
	}
}
