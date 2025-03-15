package taskpool

import (
	"log"
)

// Task represents a unit of work for the pool
type Task func()

// Pool manages a fixed number of worker goroutines
type Pool struct {
	tasks chan Task
	done  chan struct{}
}

// NewPool initializes a new worker pool
func NewPool(workerCount int) *Pool {
	p := &Pool{
		tasks: make(chan Task),
		done:  make(chan struct{}),
	}

	for i := 0; i < workerCount; i++ {
		go p.worker(i) // Start workers
	}
	return p
}

// Submit adds a new task to the pool
func (p *Pool) Submit(task Task) {
	select {
	case p.tasks <- task:
	case <-p.done:
		log.Println("Task rejected: pool is shutting down")
	}
}

// Shutdown stops all workers
func (p *Pool) Shutdown() {
	close(p.done)  // Stop accepting new tasks
	close(p.tasks) // Terminate existing workers
}

func (p *Pool) worker(id int) {
	for task := range p.tasks {
		task()
	}
	log.Printf("Worker %d shutting down", id)
}
