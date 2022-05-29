package main

import (
	"github.com/samber/lo"
	"log"
	"time"
)

func Rx(inputCh chan int) chan int {
	return rxPipe3(
		RxFiler(func(t int) bool {
			return t > 2
		}),
		RxConcatMap(func(t int) *Observable[int] {
			log.Printf("new value from filter %v", t)
			inputCh := make(chan int)
			go func() {
				for range lo.Range(100) {
					time.Sleep(time.Millisecond * 10)
					inputCh <- 1000 + t
				}
				close(inputCh)
			}()
			return New(inputCh)
		}),
		RxMap(func(t int) int {
			return t * 10
		}),
	)(inputCh)
}
func main() {
	inputCh := make(chan int)
	go func() {
		for idx, _ := range lo.Range(1000) {
			inputCh <- idx
		}
	}()
	stream := Rx(inputCh)
	go func() {
		for v := range stream {
			println(v)
		}
	}()

	time.Sleep(time.Second * 100)
}
