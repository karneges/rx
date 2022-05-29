package main

import (
	"sync"
)

type Subscription[T any] struct {
	SubscriptionCh chan T
	id             int
	Unsubscribe    func()
}
type Observer[T any] interface {
	Subscribe() chan T
}
type Observable[T any] struct {
	outCh              chan T
	subscribersMap     map[int]*Subscription[T]
	mut                sync.RWMutex
	countOfSubscribers int
}

func New[T any](inputChannel chan T) *Observable[T] {
	outCh := make(chan T)
	subscribersMap := make(map[int]*Subscription[T])
	observable := &Observable[T]{
		outCh:          outCh,
		mut:            sync.RWMutex{},
		subscribersMap: subscribersMap,
	}
	go func() {
		for {
			inputChannelValue, hasMore := <-inputChannel
			if hasMore {
				for _, value := range observable.subscribersMap {
					observable.mut.Lock()
					if value != nil {
						value.SubscriptionCh <- inputChannelValue
					}
					observable.mut.Unlock()
				}
			} else {
				for _, value := range observable.subscribersMap {
					if value != nil {
						value.Unsubscribe()
					}
				}
			}

		}
	}()
	return observable

}

func (o *Observable[T]) unsubscribe(id int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	close(o.subscribersMap[id].SubscriptionCh)

	delete(o.subscribersMap, id)
	o.countOfSubscribers -= 1
}
func (o *Observable[T]) Subscribe() *Subscription[T] {
	o.countOfSubscribers += 1
	println(o.countOfSubscribers)
	newSubscription := &Subscription[T]{
		SubscriptionCh: make(chan T),
		id:             o.countOfSubscribers,
		Unsubscribe: func() {
			o.unsubscribe(o.countOfSubscribers)
		},
	}
	o.subscribersMap[newSubscription.id] = newSubscription
	return newSubscription
}
