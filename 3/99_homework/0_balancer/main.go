package main

import (
	"fmt"
	"time"
)

const timeout = 100 * time.Millisecond

type RoundRobinBalancer struct {
	pool  []int
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

/*Init - инициализирует собственно балансер -
представьте что устанавливает соединения с указанным количеством серверов*/
func (r *RoundRobinBalancer) Init(n int) {
	*r = RoundRobinBalancer{
		pool: make([]int, n),
		size: 4,
		// Канал задач - буферизированный, чтобы основная программа не блокировалась при постановке задач
		tasks: make(chan Task, 128),
		// Канал kill для убийства "лишних воркеров"
		kill: make(chan struct{}),
	}
	// fmt.Println(r.uu)
	return
}

// GiveStat - даёт статистику, сколько запросов пришло на каждый из серверов.
func (r *RoundRobinBalancer) GiveStat() (statist []int) {
	return r.pool
}

// GiveNode - эта функция вызывается, когда пришел запрос. мы получаем номер сервера, на который идти.
func (r *RoundRobinBalancer) GiveNode() (n int) {
	n = r.size
	fmt.Println(r.size)

	return
}

type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

// Скроем внутреннее усройство за конструктором, пользователь может влиять только на размер пула
func NewPool(size int) *Pool {
	pool := &Pool{
		// Канал задач - буферизированный, чтобы основная программа не блокировалась при постановке задач
		tasks: make(chan Task, 128),
		// Канал kill для убийства "лишних воркеров"
		kill: make(chan struct{}),
	}
	// Вызовем метод resize, чтобы установить соответствующий размер пула
	pool.Resize(size)
	return pool
}

// Жизненный цикл воркера
func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		// Если есть задача, то ее нужно обработать
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
			// Если пришел сигнал умирать, выходим
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) Resize(n int) {
	// Захватывам лок, чтобы избежать одновременного изменения состояния
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

func (p *Pool) Close() {
	close(p.tasks)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Exec(task Task) {
	p.tasks <- task
}
