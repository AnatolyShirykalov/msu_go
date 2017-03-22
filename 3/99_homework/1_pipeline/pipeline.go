package pipeline

import "sync"

type job func(in, out chan interface{})

func Pipe(funcs ...job) {
	var wg sync.WaitGroup
	chs := make([]chan interface{}, len(funcs))
	wg.Add(len(funcs))
	for i := range funcs {
		chs[i] = make(chan interface{})
		go func(k int) {
			if k == 0 {
				funcs[0](chs[0], chs[0])
			} else {
				funcs[k](chs[k-1], chs[k])
			}
			close(chs[k])
			wg.Done()
		}(i)
	}
	wg.Wait()
	return
}
