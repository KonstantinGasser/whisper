package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/KonstantinGasser/whisper"
)

func main() {

	broker := whisper.NewBroker()

	broker.NewTopic("topic-1")
	broker.NewTopic("topic-2")

	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer fmt.Printf("consumer[%d] no more message - bye!\n", id)
			defer wg.Done()

			consumer, err := broker.Subscribe("topic-1")
			if err != nil {
				panic(err)
			}

			for {
				msg, ok := consumer.Consume()
				if !ok {
					break
				}
				fmt.Printf("Consumer[topic-1]: %v\n", msg.Data)
			}
		}(i, &wg)
	}

	go func() {
		for i := 0; i < 10000; i++ {
			if err := broker.Publish("topic-1", whisper.Message{
				Data:      fmt.Sprintf("Data %d", i),
				Timestamp: time.Now().Unix(),
			}); err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	time.Sleep(time.Second * 5)
	fmt.Println("Stopping broker")
	broker.Stop()
	fmt.Println("Closed")
}
