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

// Option 功能选项
type Option func(*Pool)

func WithBlock(block bool) Option {
	return func(pool *Pool) {
		pool.block = block
	}
}

func WithPreAllocWorkers(preAlloc bool) Option {
	return func(pool *Pool) {
		pool.preAlloc = preAlloc
	}
}

type Pool struct {
	capacity int

	active chan struct{} // 计数信号量，用于控制 worker 的数量

	tasks chan Task

	wg sync.WaitGroup

	quit chan struct{}

	// 是否在创建pool的时候就预创建workers，默认值为：false
	preAlloc bool

	// 当pool满的情况下，新的Schedule调用是否阻塞当前goroutine。默认值：true
	// 如果block = false，则Schedule返回ErrNoWorkerAvailInPool
	block bool
}

// New 初始化线程池
func New(capacity int, opts ...Option) *Pool {
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

	for _, opt := range opts {
		opt(p)
	}

	fmt.Printf("workerpool start\n")

	// 预创建 worker
	if p.preAlloc {
		for i := 0; i < p.capacity; i++ {
			p.newWorker(i + 1)
			p.active <- struct{}{}
		}
	}

	go p.run()

	return p
}

// ErrWorkerPoolFreed 哨兵错误
var (
	ErrWorkerPoolFreed    = errors.New("线程池已经释放")
	ErrNoIdleWorkerInPool = errors.New("没有空闲的worker处理任务")
)

// Schedule 提交任务
// TODO：在没有达到最大 worker 数时，可以考虑创建，而不是直接阻塞
// TODO：调度的时候阻塞，可以考虑将新增的 task 放入到队列中，而不是直接报错
func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	default:
		if p.block {
			p.tasks <- t
			return nil
		}
		return ErrNoIdleWorkerInPool
	}
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait() // 等待所有 worker 退出
	fmt.Printf("线程池退出成功")
}

func (p *Pool) run() {
	idx := len(p.active)

	// 根据 task 创建 worker
	if !p.preAlloc {
	loop:
		for task := range p.tasks {
			// 将任务塞回去
			go func() {
				p.tasks <- task
			}()
			select {
			case <-p.quit:
				return
			case p.active <- struct{}{}:
				idx++
				p.newWorker(idx)
			default:
				break loop
			}
		}
	}
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
