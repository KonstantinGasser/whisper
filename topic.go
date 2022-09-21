package ares

import (
	"sync"
	"time"
)

func WithTimeout(timeout time.Duration, onTimeout func()) func(*Topic) {
	return func(t *Topic) {

	}
}

func WithBuffer(buffer uint64) func(*Topic) {
	return func(t *Topic) {
	}
}

type Topic struct {
	sync.RWMutex
	consumer []chan interface{}
}

func newTopic(label string, opts ...func(*Topic)) *Topic {
	return &Topic{
		consumer: make([]chan interface{}, 0),
	}
}

func (t *Topic) publish(event interface{}) {

	for _, cons := range t.consumer {
		cons <- event
	}
}

func (t *Topic) subscribe() Consumer {

	evtChan := make(chan interface{})

	c := consumer{
		stream: evtChan,
	}

	t.RLock()
	defer t.RUnlock()
	t.consumer = append(t.consumer, evtChan)

	return c
}
