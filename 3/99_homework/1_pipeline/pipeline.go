package pipeline

import (
	//"time"
	"sync"
)

type job func(in, out chan interface{})

func Pipe(funcs ...job) {
	var wg sync.WaitGroup
	wg.Add(1)
	chs := make([]chan interface{}, len(funcs))
	for i, fun := range funcs {
		chs[i] = make(chan interface{})
		if i == 0 {
			go func() {
				funcs[0](chs[0], chs[0])
				wg.Done()
			}()
			continue
		}
		go fun(chs[i-1], chs[i])
	}
	wg.Wait()
	return
}
