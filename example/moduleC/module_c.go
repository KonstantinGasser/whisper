package modulec

import (
	"context"
	"fmt"

	"github.com/KonstantinGasser/whisper"
)

type Module struct {
	consumer whisper.Consumer
}

func New(broker whisper.Broker) *Module {

	cns, err := broker.Subscribe("topic-1")

	if err != nil {
		panic(err)
	}
	return &Module{
		consumer: cns,
	}
}

func (m Module) Start(ctx context.Context) {

	for {
		msg, ok := m.consumer.Consume()
		if !ok {
			return
		}

		fmt.Printf("Module-C: Message from %s: Data: %v\n", msg.Topic(), msg.Data)
	}
}
