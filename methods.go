package main

func rxPipe2[P1 any, C1 any, C2 any](f1 func(chan P1) chan C1, f2 func(chan C1) chan C2) func(chan P1) chan C2 {
	return func(p1 chan P1) chan C2 {
		return f2(f1(p1))
	}
}
func rxPipe3[P1 any, C1 any, C2 any, C3 any](
	f1 func(chan P1) chan C1,
	f2 func(chan C1) chan C2,
	f3 func(chan C2) chan C3,
) func(chan P1) chan C3 {
	return func(p1 chan P1) chan C3 {
		return f3(f2(f1(p1)))
	}
}

func RxFiler[T any](f func(v T) bool) func(chan T) chan T {
	outCh := make(chan T)
	return func(inputCh chan T) chan T {
		go func() {
			for chanValue := range inputCh {
				if f(chanValue) {
					outCh <- chanValue
				}
			}
		}()
		return outCh

	}
}
func RxMap[T any, R any](f func(T) R) func(chan T) chan R {
	outCh := make(chan R)
	return func(inputCh chan T) chan R {
		go func() {
			for chanValue := range inputCh {
				outCh <- f(chanValue)
			}
		}()
		return outCh

	}
}

func RxFirst[T any](f func(v T) bool) func(chan T) chan T {
	outCh := make(chan T)
	return func(inputCh chan T) chan T {
		go func() {
			for chanValue := range inputCh {
				if f(chanValue) {
					outCh <- chanValue
					return
				}
			}
		}()
		return outCh

	}
}

func RxSwitchMap[T any, E any](f func(T) *Observable[E]) func(chan T) chan E {
	return func(inputCh chan T) chan E {
		var subscription *Subscription[E]
		outCh := make(chan E)

		go func() {
			for inputValue := range inputCh {
				if subscription != nil {
					subscription.Unsubscribe()
				}
				subscription = f(inputValue).Subscribe()
				go func() {
					for subscriptionValue := range subscription.SubscriptionCh {
						outCh <- subscriptionValue
					}
				}()

			}
		}()
		return outCh
	}
}

func RxMergeMap[T any, E any](f func(T) *Observable[E]) func(chan T) chan E {
	return func(inputCh chan T) chan E {
		var subscription *Subscription[E]
		outCh := make(chan E)

		go func() {
			for inputValue := range inputCh {
				subscription = f(inputValue).Subscribe()
				go func() {
					for subscriptionValue := range subscription.SubscriptionCh {
						outCh <- subscriptionValue
					}
				}()

			}
		}()
		return outCh
	}
}

func RxConcatMap[T any, E any](f func(T) *Observable[E]) func(chan T) chan E {
	return func(inputCh chan T) chan E {
		var subscription *Subscription[E]
		outCh := make(chan E)

		go func() {
			for inputValue := range inputCh {
				subscription = f(inputValue).Subscribe()
				for subscriptionValue := range subscription.SubscriptionCh {
					outCh <- subscriptionValue
				}
			}
		}()
		return outCh
	}
}
