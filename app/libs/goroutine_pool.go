package libs

import (
	"log"
	"sync"
)

// 一个简单的线程池
type (
	Task struct {
		Id   int
		Task func()
	}

	GoroutinePool struct {
		queue     chan *Task
		wg        sync.WaitGroup
		workerNum int
	}
)

func NewGoroutinePool(num int) *GoroutinePool {
	goroutinePool := &GoroutinePool{
		queue:     make(chan *Task, num),
		wg:        sync.WaitGroup{},
		workerNum: num,
	}

	return goroutinePool
}

func (g *GoroutinePool) Wait() {
	close(g.queue)
	g.wg.Wait()
}

func (g *GoroutinePool) worker() {
	g.wg.Add(1)
	defer g.wg.Done()

	for v := range g.queue {
		g.execute(v)
	}
}

func (g *GoroutinePool) execute(task *Task) {
	// 防止panic导致整个worker崩溃
	defer func() {
		if rcv := recover(); rcv != nil {
			log.Println("%v", rcv)
		}
	}()

	task.Task()
}

func (g *GoroutinePool) AddTask(task *Task) {
	g.queue <- task
}

// num 为worker的数量
func (g *GoroutinePool) Run() {
	for i := 0; i < g.workerNum; i++ {
		go g.worker()
	}
}
