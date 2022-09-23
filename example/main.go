package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/KonstantinGasser/whisper"
	modulea "github.com/KonstantinGasser/whisper/example/moduleA"
	moduleb "github.com/KonstantinGasser/whisper/example/moduleB"
	modulec "github.com/KonstantinGasser/whisper/example/moduleC"
)

func main() {

	broker := whisper.NewBroker()

	broker.NewTopic("topic-1")
	broker.NewTopic("topic-2")

	mA := modulea.New(broker)
	mB := moduleb.New(broker)
	mC := modulec.New(broker)

	ctx := context.Background()
	go mA.Start(ctx)
	go mB.Start(ctx)
	go mC.Start(ctx)

	for i := 0; i < 10; i++ {

		topic := "topic-1"
		if rand.Intn(2) == 1 {
			topic = "topic-2"
		}

		if err := broker.Publish(topic, whisper.Message{
			Timestamp: time.Now().Unix(),
			Data:      i,
		}); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	}
}
