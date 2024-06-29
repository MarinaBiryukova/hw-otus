package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorsCounter struct {
	mu    sync.RWMutex
	m     int
	count int
}

func (c *errorsCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *errorsCounter) Exceeded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.m == 0 {
		return true
	}

	if c.m < 0 {
		return false
	}

	return c.count >= c.m
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, len(tasks))
	for i := range tasks {
		tasksCh <- tasks[i]
	}
	close(tasksCh)

	errCounter := errorsCounter{m: m}
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				if errCounter.Exceeded() {
					return
				}

				task, ok := <-tasksCh
				if !ok {
					return
				}

				err := task()
				if err != nil {
					errCounter.Increment()
				}
			}
		}()
	}

	wg.Wait()

	if errCounter.Exceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
