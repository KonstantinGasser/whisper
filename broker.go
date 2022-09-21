package whisper

import (
	"fmt"
	"sync"
)

// Broker is an in-memory message broker orchestrating
// Messages via topics. Each topic is identified by a unique route.
// Producer can publish messages to the broker and Consumers can receive
// Messages from topics. A Messages is distributed to all connected consumers
// however, a topic does not store Messages which is why the Broker is a "fire and forget"
// type of Message Broker. Moreover, at the moment a message published to a topic can be blocking
// if a consumer of the topic is taking more time to process a message. As such it is recommended
// that each result from Consumer.Consume is handled in its own goroutine to avoid waiting times
type Broker interface {
	// NewTopic registers a new topic in the Broker but returns
	// an error if a topic with the same route already exists
	NewTopic(route string, opts ...func(*Topic)) error

	// Subscribe allows consumers to start receiving messages from
	// the requested topic. The returned Consumer has a function
	// Consume which is a blocking call to receive messages from the topic
	Subscribe(route string) (Consumer, error)

	// Publish allows to push a Message to a specific topic which
	// can be consumed by all connected Consumers
	Publish(topic string, msg Message) error

	// Stop propagates down to each topic and closes the topic.
	// As such all consumers will return false in their next call of Consume
	Stop()
}

// implements the Broker interface
type broker struct {
	sync.RWMutex
	topics map[string]*Topic
	closed bool
}

func NewBroker() Broker {
	return &broker{
		topics: make(map[string]*Topic),
		closed: false,
	}
}

func (b *broker) NewTopic(route string, opts ...func(*Topic)) error {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.topics[route]; ok {
		return fmt.Errorf("topic %q already exists", route)
	}

	t := newTopic(route, opts...)
	b.topics[route] = t

	return nil
}

func (b *broker) Subscribe(route string) (Consumer, error) {
	b.Lock()
	defer b.Unlock()

	t, ok := b.topics[route]
	if !ok {
		return nil, fmt.Errorf("topic %q does not exist", route)
	}

	return t.subscribe(), nil
}

// Publish is a blocking operation if the receiving topic
// and its consumers are blocked
func (b *broker) Publish(route string, msg Message) error {
	b.Lock()
	defer b.Unlock()

	if b.closed {
		return fmt.Errorf("publish on closed broker topics")
	}
	topic, ok := b.topics[route]
	if !ok {
		return fmt.Errorf("topic %q does not exist", route)
	}

	topic.publish(msg)
	return nil
}

func (b *broker) Stop() {
	b.Lock()
	defer b.Unlock()

	b.closed = true
	for _, topic := range b.topics {
		topic.stop()
	}
}
