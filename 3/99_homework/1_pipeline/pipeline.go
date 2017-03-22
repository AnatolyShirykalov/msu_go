package pipeline

import (
	//"time"
	"sync"
)

type job func(in, out chan interface{})

func Pipe(funcs ...job) {
	var wg sync.WaitGroup
	wg.Add(len(funcs))
	chs := make([]chan interface{}, len(funcs))
	for i, _ := range funcs {
		chs[i] = make(chan interface{})
		go func(j int) {
			k := j - 1
			if j == 0 {
				k = 0
			}
			funcs[j](chs[k], chs[j])
			close(chs[j])
			wg.Done()
		}(i)
	}
	wg.Wait()
	return
}
