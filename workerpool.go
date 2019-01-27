package workerpool

import (
	"sync"
	"time"
)

// WorkerPool can manage a collection of workers, each with their own goroutine.
type WorkerPool interface {
	Spawn(Fn, Argv) bool
	Start()
	Stop()
}

// New creates a new WorkerPool
func New(size int, maxIdleWorkerDuration time.Duration) WorkerPool {
	return &implWorkerPool{
		workers:               new(workerlist),
		maxSize:               size,
		maxIdleWorkerDuration: maxIdleWorkerDuration,
	}
}

type implWorkerPool struct {
	lock sync.Mutex

	workerCount           int
	maxSize               int
	mustStop              bool
	vworkerPool           sync.Pool
	workers               *workerlist
	maxIdleWorkerDuration time.Duration
	stopCh                chan struct{}
}

func (pool *implWorkerPool) Spawn(f Fn, a Argv) bool {
	w := pool.getWorker()

	if w == nil {
		return false
	}

	w.taskCh <- task{f, a}

	return true
}

func (pool *implWorkerPool) Start() {
	if pool.stopCh != nil {
		panic("BUG: workerPool already started")
	}

	pool.stopCh = make(chan struct{})
	stopCh := pool.stopCh
	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				time.Sleep(pool.maxIdleWorkerDuration)
			}
		}
	}()
}

func (pool *implWorkerPool) Stop() {
	if pool.stopCh == nil {
		panic("BUG: workerPool wasn't started")
	}

	close(pool.stopCh)
	pool.stopCh = nil

	pool.lock.Lock()
	var ele *element
	for ele = pool.workers.Front(); ele != nil; ele = ele.next {
		ele.value.taskCh <- task{fn: nil}
	}
	pool.workers.root.prev = nil
	pool.workers.root.next = nil
	pool.lock.Unlock()
}

func (pool *implWorkerPool) getWorker() *worker {
	createWorker := false
	var w *worker

	pool.lock.Lock()
	w, ok := pool.workers.PopBack()
	if !ok {
		if pool.workerCount < pool.maxSize {
			createWorker = true
			pool.workerCount++
		}
	}
	pool.lock.Unlock()

	if w == nil {
		if !createWorker {
			return nil
		}

		vworker := pool.vworkerPool.Get()
		if vworker == nil {
			vworker = &worker{
				taskCh: make(chan task),
			}
		}
		w = vworker.(*worker)
		go func() {
			pool.workerFunc(w)
			pool.vworkerPool.Put(w)
		}()
	}

	return w
}

func (pool *implWorkerPool) clean() {
	currentTime := time.Now()
	maxIdleWorkerDuration := pool.maxIdleWorkerDuration

	pool.lock.Lock()

	var ele *element
	for ele = pool.workers.Front(); ele != nil; ele = ele.next {
		w := ele.value
		if currentTime.Sub(w.lastUseTime) < maxIdleWorkerDuration {
			pool.workers.ResetFront(ele)
			break
		}

		w.taskCh <- task{fn: nil}
	}
}

func (pool *implWorkerPool) release(w *worker) bool {
	w.lastUseTime = time.Now()

	pool.lock.Lock()
	if pool.mustStop {
		pool.lock.Unlock()
		return false
	}

	pool.workers.PushBack(w)
	pool.lock.Unlock()

	return true
}

func (pool *implWorkerPool) workerFunc(w *worker) {
	for t := range w.taskCh {
		if t.fn == nil {
			break
		}

		t.fn(t.argv)
		if !pool.release(w) {
			break
		}
	}

	pool.lock.Lock()
	pool.workerCount--
	pool.lock.Unlock()
}

// Argv is arguments value
type Argv interface{}

// Fn is a function
type Fn func(argv Argv)

type task struct {
	fn   Fn
	argv Argv
}

type worker struct {
	taskCh      chan task
	lastUseTime time.Time
}
