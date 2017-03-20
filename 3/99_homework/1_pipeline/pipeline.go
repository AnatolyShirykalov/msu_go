package pipeline

import (
	"time"
)

type job func(in, out chan interface{})

func Pipe(funcs ...job) {
	slChan := make([]chan interface{}, len(funcs))
	for i, fun := range funcs {
		slChan[i] = make(chan interface{})
		if i == 0 {
			go fun(slChan[0], slChan[0])
			continue
		}
		go fun(slChan[i-1], slChan[i])
	}
	time.Sleep(1 * time.Microsecond)
}
