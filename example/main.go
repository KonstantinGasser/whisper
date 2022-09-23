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

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			// defer fmt.Printf("consumer[%d] no more message - bye!\n", id)
			defer wg.Done()

			consumer, err := broker.Group("topic-1", "topic-2")
			if err != nil {
				panic(err)
			}

			for {
				msg, ok := consumer.Consume()
				if !ok {
					break
				}
				fmt.Printf("Consumer[%s]: %v\n", msg.Topic, msg.Data)
			}
		}(i, &wg)
	}

	for i := 0; i < 5; i++ {
		if err := broker.Publish("topic-1", whisper.Message{
			Data:      fmt.Sprintf("Data %d", i),
			Timestamp: time.Now().Unix(),
		}); err != nil {
			fmt.Println(err)
			return
		}

		if err := broker.Publish("topic-2", whisper.Message{
			Data:      fmt.Sprintf("Data %d", i),
			Timestamp: time.Now().Unix(),
		}); err != nil {
			fmt.Println(err)
			return
		}
	}

	time.Sleep(time.Second * 5)
	fmt.Println("Stopping broker")
	broker.Stop()
	fmt.Println("Closed")
}
