package ares

import (
	"fmt"
	"sync"
)

type Broker interface {
	NewTopic(label string, opts ...func(*Topic)) error
	Subscribe(label string) (Consumer, error)

	Publish(topic string, data interface{}) error
}

type broker struct {
	sync.RWMutex
	topics map[string]*Topic
}

func NewBroker() Broker {
	return &broker{
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

func (b *broker) Publish(label string, data interface{}) error {
	b.Lock()
	defer b.Unlock()

	t, ok := b.topics[label]
	if !ok {
		return fmt.Errorf("topic %q does not exist", label)
	}

	t.publish(data)
	return nil
}
