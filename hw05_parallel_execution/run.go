package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrEmptyTasks          = errors.New("empty task list")
)

type Task func() error

// Run Tasks with channels.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrEmptyTasks
	}

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	workerMaxCount := n
	// если кол-во задач меньше кол-ва воркеров, ограничиваем кол-во воркеров
	if len(tasks) < n {
		workerMaxCount = len(tasks)
	}

	tasksCh := make(chan Task)       // канал задач.
	stopErrCh := make(chan struct{}) // сигнальный канал.
	errorCh := make(chan error)      // канал ошибок.
	wg := sync.WaitGroup{}
	wgErr := sync.WaitGroup{}

	defer func() {
		close(tasksCh)
		close(stopErrCh)
		wg.Wait()
		close(errorCh)
		wgErr.Wait()
	}()

	// запускаем воркеры
	for i := 0; i < workerMaxCount; i++ {
		wg.Add(1)
		go doWork(&wg, tasksCh, errorCh)
	}

	wgErr.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		errCount := 0 // текущее кол-во ошибок
		for range errorCh {
			errCount++
			if errCount == m {
				stopErrCh <- struct{}{}
			}
		}
	}(&wgErr)

	for _, task := range tasks {
		select {
		case <-stopErrCh:
			return ErrErrorsLimitExceeded
		case tasksCh <- task:
		}
	}

	return nil
}

func doWork(wg *sync.WaitGroup, tasks <-chan Task, errs chan<- error) {
	defer func() {
		wg.Done()
	}()

	for task := range tasks {
		err := task()
		if err != nil {
			errs <- err
		}
	}
}
