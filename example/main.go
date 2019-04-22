package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ruizu/go-broadcast/broadcast"
)

func main() {
	// initialize a new broadcast instance
	b := broadcast.New(5)

	// call a function that subscribes to the broadcast
	// in this example, the function will start randomly
	go doSomething(1, b)
	go doSomething(2, b)
	go doSomething(3, b)

	// publish messages to the broadcast
	for i := 0; i < 100; i++ {
		fmt.Println("Publishing: ", i)
		b.Publish(i)
		time.Sleep(500 * time.Millisecond)
	}
}

func doSomething(id int, b broadcast.Broadcast) {
	// randomly sleep
	time.Sleep(time.Duration(rand.Intn(id*3)) * time.Second)

	// start subscribing to the broadcast
	ch := b.Subscribe()
	defer func() {
		// unsubscribe the channel after the function is done
		b.Unsubscribe(ch)
	}()

	// receive data from the broadcast
	for {
		msg := <-ch
		fmt.Println("Worker ", id, ": ", msg)
	}
}
