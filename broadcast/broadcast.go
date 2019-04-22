package broadcast

import (
	"sync"
)

// A Broadcast interface is used to descibe the functions available.
type Broadcast interface {
	// Publish a message to subscribing channels.
	Publish(interface{})
	// Subscribe returns a new receive-only channel that subscribes to the broadcast.
	Subscribe() <-chan interface{}
	// Unsubscribe closes the subscribed channel and remove it from the subscriber list.
	Unsubscribe(<-chan interface{})
	// Close shuts down the broadcast and closes all subscriber channel.
	Close()
}

type broadcast struct {
	mu          sync.Mutex
	subscribers map[<-chan interface{}]chan interface{}
	messages    chan interface{}
}

// New initialize a new broadcast instance.
func New(n int) Broadcast {
	b := &broadcast{
		subscribers: make(map[<-chan interface{}]chan interface{}),
		messages:    make(chan interface{}, n),
	}
	go b.publishMessages()
	return b
}

func (b *broadcast) Publish(msg interface{}) {
	b.messages <- msg
}

func (b *broadcast) Subscribe() <-chan interface{} {
	newCh := make(chan interface{})
	b.mu.Lock()
	b.subscribers[newCh] = newCh
	b.mu.Unlock()
	return newCh
}

func (b *broadcast) Unsubscribe(ch <-chan interface{}) {
	b.mu.Lock()
	rwCh, ok := b.subscribers[ch]
	if ok {
		delete(b.subscribers, ch)
		close(rwCh)
	}
	b.mu.Unlock()
}

func (b *broadcast) Close() {
	b.mu.Lock()
	close(b.messages)
	for ch, rwCh := range b.subscribers {
		delete(b.subscribers, ch)
		close(rwCh)
	}
	b.mu.Unlock()
}

func (b *broadcast) publishMessages() {
	for {
		msg := <-b.messages
		for _, rwCh := range b.subscribers {
			rwCh <- msg
		}
	}
}
