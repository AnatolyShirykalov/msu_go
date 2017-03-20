package main

//есть некоторое количество серверов, нагрузка на которых распределяется методом "по кругу"
//есть балансер, который должен распределить запросы равномерно
type RoundRobinBalancer struct {
	stat    []int
	next    int
	tasks   chan int
	directs chan int
}

//Init - инициализирует собственно балансер - представьте что устанавливает соединения с указанным колчиеством серверов.
func (r *RoundRobinBalancer) Init(n int) {
	r.stat = make([]int, n)
	r.next = 0
	r.tasks = make(chan int, 1)
	r.directs = make(chan int, 1)

	go func() {
		for range r.tasks {

			r.stat[r.next]++
			r.next = (r.next + 1) % len(r.stat)

			r.directs <- r.next

		}
	}()
}

//GiveStat - даёт статистику, сколько запросов пришло на каждый из серверов.
func (r *RoundRobinBalancer) GiveStat() []int {
	return r.stat
}

//GiveNode - эта функция фвзывается, когда пришел запрос. мы получаем номер сервера, на который идти.
func (r *RoundRobinBalancer) GiveNode() int {

	r.tasks <- 1
	return <-r.directs
}
