package types

import (
	"net/http"
	"sync"
	"time"
)

type Bus struct {
	Subscriptions map[string]Subscription
	Mutex         *sync.RWMutex
}

type Subscription struct {
	Data    chan *http.Response
	Created *time.Time
}

func NewBus() *Bus {
	return &Bus{
		Subscriptions: make(map[string]Subscription),
		Mutex:         &sync.RWMutex{},
	}
}

func (b *Bus) SubscriptionList() []string {
	keys := []string{}

	b.Mutex.RLock()

	for key := range b.Subscriptions {
		keys = append(keys, key)
	}

	b.Mutex.RUnlock()
	return keys
}

func (b *Bus) Send(id string, res *http.Response) {
	var ok bool

	b.Mutex.RLock()
	defer b.Mutex.RUnlock()
	
	_, ok = b.Subscriptions[id]

	if !ok {
		return
	}

	b.Subscriptions[id].Data <- res
}

func (b *Bus) Subscribe(id string) *Subscription {
	now := time.Now()
	sub := Subscription{
		Data:    make(chan *http.Response),
		Created: &now,
	}

	b.Mutex.Lock()
	b.Subscriptions[id] = sub
	b.Mutex.Unlock()

	return &sub
}

func (b *Bus) Expired(id string, after time.Duration) bool {

	b.Mutex.RLock()

	sub, ok := b.Subscriptions[id]

	b.Mutex.RUnlock()

	if ok && sub.Created != nil {
		return (*sub.Created).Add(after).After(time.Now())
	}

	return false
}

func (b *Bus) Unsubscribe(id string) {

	b.Mutex.RLock()

	sub, ok := b.Subscriptions[id]

	b.Mutex.RUnlock()

	if ok {
		close(sub.Data)

		b.Mutex.Lock()

		delete(b.Subscriptions, id)
		b.Mutex.Unlock()
	}
}
