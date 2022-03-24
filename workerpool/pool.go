package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

const (
	defaultCapacity = 5
	maxCapacity     = 20
)

type Task func()

type Pool struct {
	capacity int

	active chan struct{} // 计数信号量，用于控制 worker 的数量

	tasks chan Task

	wg sync.WaitGroup

	quit chan struct{}
}

// New 初始化线程池
func New(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}
	p := &Pool{
		capacity: capacity,
		active:   make(chan struct{}, capacity),
		tasks:    make(chan Task),
		wg:       sync.WaitGroup{},
		quit:     make(chan struct{}),
	}

	fmt.Printf("workerpool start\n")

	go p.run()

	return p
}

// ErrWorkerPoolFreed 哨兵错误
var ErrWorkerPoolFreed = errors.New("线程池已经释放")

// Schedule 提交任务
func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	}
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait() // 等待所有 worker 退出
	fmt.Printf("线程池退出成功")
}

func (p *Pool) run() {
	idx := 0

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
			}

			<-p.active  // goroutine 退出时，计数信号量-1
			p.wg.Done() // 减少 waitGroup 等待 goroutine 的数量
		}()
		fmt.Printf("worker[%03d]: start\n", i)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", i)
				return
			case t := <-p.tasks: // 不断获取 task 并执行
				fmt.Printf("worker[%03d]: receive a task\n", i)
				t()
			}
		}
	}()
}
