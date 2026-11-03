package main

import (
	"fmt"
	"gomoney-assessment/broker"
)

func main() {
	// Create a new broker for string messages
	b := broker.NewBroker[string]()

	// Create a publisher for the broker
	publisher := broker.NewPublisher(b)

	// Publish a range of messages
	publishRange(publisher, 1, 5)

	sub1 := broker.NewSubscriber[string]("Subscriber 1")
	sub1.Subscribe(b)

	publishRange(publisher, 6, 10)

	sub2 := broker.NewSubscriber[string]("Subscriber 2")
	sub2.Subscribe(b)

	printOut(sub1)
	printOut(sub2)
}

func publishRange(publisher *broker.Publisher[string], start, end int) {
	for i := start; i <= end; i++ {
		publisher.Publish(fmt.Sprintf("Message %d", i))
	}
}

func printOut(sub *broker.Subscriber[string]) {
	for {
		message, ok := sub.Next()
		if !ok {
			break
		}
		fmt.Printf("%s received: %s\n", sub.Name(), message.Payload)
	}
}
