package pipeline


import (
  "time"
)
type job func(in, out chan interface{})

func Pipe(funcs ...job) {
        chs := make([]chan interface{},len(funcs))
        for i, fun := range funcs {
                chs[i] = make(chan interface{})
                if i==0 {
                        go fun(chs[0], chs[0])
                        continue
                }
                if i==len(funcs) - 1 {
                        fun(chs[i-1], chs[i])
                        continue
                }
                go fun(chs[i-1], chs[i])
        }
                                                                                                                                                                                                                                                                 time.Sleep(200 * time.Millisecond)
	return
}
