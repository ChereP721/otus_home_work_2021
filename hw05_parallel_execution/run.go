package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrCannotRunWorkers = errors.New("cannot run workers, negative or zero number of workers")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return ErrCannotRunWorkers
	}

	termCh := make(chan struct{}, 1)
	taskCh := genTasksCh(tasks, termCh)
	ec := newErrorCounter(m, termCh)

	sg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		sg.Add(1)

		go func() {
			defer sg.Done()

			for {
				if ec.isFailed() {
					return
				}

				select {
				case task, ok := <-taskCh:
					if !ok {
						return
					}
					err := task()
					if err != nil {
						ec.inc()
					}
				default:
				}
			}
		}()
	}

	sg.Wait()
	/*
		тест где м больше н и наоборот
		где кол-во ошибок меньше допустимого числа
			тест без слипа - через библиотеки работы со временем из лекции по тестам 2
	*/

	if ec.isFailed() {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func genTasksCh(tasks []Task, termCh <-chan struct{}) <-chan Task {
	taskCh := make(chan Task)

	go func() {
		for _, task := range tasks {
			select {
			case <-termCh:
				close(taskCh)
				return
			case taskCh <- task:
			}
		}
		close(taskCh)
	}()

	return taskCh
}
