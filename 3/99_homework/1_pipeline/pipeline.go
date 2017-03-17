package pipeline

type job func(in, out chan interface{})

func Pipe(funcs ...job) {
        cin := make(chan interface{},4)
        cout := make(chan interface{},4)
        for _, n := []interface{}{0, 30, 60, 90} {
                cin <- n
        }
        for fn := range funcs {
                go fn(cin, cout)
        }
	return
}
