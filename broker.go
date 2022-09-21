package whisper

import (
	"fmt"
	"sync"
)

// Broker is an in-memory message broker orchestrating
// Messages through topics. Each topic is identified by a unique route.
// Producer can publish messages to the broker and Consumers can receive
// Messages from topics. A Messages is distributed to all connected consumers
// however, a topic does not store Messages which is why the Broker is a "fire and forget"
// type of Message Broker.
type Broker interface {
	// NewTopic registers a new topic in the Broker but returns
	// an error if a topic with the same route already exists
	NewTopic(label string, opts ...func(*Topic)) error

	// Subscribe allows consumers to start receiving messages from
	// the requested topic
	Subscribe(label string) (Consumer, error)

	// Publish allows to push a Message to a specific topic which
	// can be consumed by all connected Consumers
	Publish(topic string, msg Message) error

	Stop()
}

// implements the Broker interface
type broker struct {
	sync.RWMutex
	queue  []Message
	topics map[string]*Topic
}

func NewBroker() Broker {
	return &broker{
		queue:  make([]Message, 0),
		topics: make(map[string]*Topic),
	}
}

func (b *broker) NewTopic(label string, opts ...func(*Topic)) error {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.topics[label]; ok {
		return fmt.Errorf("topic %q already exists", label)
	}

	t := newTopic(label, opts...)
	b.topics[label] = t

	return nil
}

func (b *broker) Subscribe(label string) (Consumer, error) {
	b.Lock()
	defer b.Unlock()

	t, ok := b.topics[label]
	if !ok {
		return nil, fmt.Errorf("topic %q does not exist", label)
	}

	return t.subscribe(), nil
}

// Publish is a blocking operation if the receiving topic
// and its consumers are blocked
func (b *broker) Publish(label string, msg Message) error {
	b.Lock()
	defer b.Unlock()

	topic, ok := b.topics[label]
	if !ok {
		return fmt.Errorf("topic %q does not exist", label)
	}

	topic.publish(msg)
	return nil
}

func (b *broker) Stop() {
	b.Lock()
	defer b.Unlock()

	for _, topic := range b.topics {
		topic.stop()
	}
}
