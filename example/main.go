package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/KonstantinGasser/ares/ares"
)

func main() {

	broker := ares.NewBroker()

	broker.NewTopic("setup.meta-data.created", ares.WithTimeout(time.Millisecond*50, nil))

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer fmt.Printf("consumer[%d] no more message - bye!", id)
			defer wg.Done()

			consumer, err := broker.Subscribe("setup.meta-data.created")
			if err != nil {
				panic(err)
			}

			for {
				data, ok := consumer.Consume()
				if !ok {
					break
				}

				fmt.Printf("Consumer[%d]: %v\n", id, data)
			}
		}(i, &wg)
	}

	for i := 0; i < 1; i++ {
		if err := broker.Publish("setup.meta-data.created", fmt.Sprintf("Event: %d", i)); err != nil {
			panic(err)
		}
	}

	wg.Wait()

}
