package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrEmptyTasks          = errors.New("empty task list")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrEmptyTasks
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var (
		wg sync.WaitGroup

		errMaxCount    = int32(m)   // максимально допустимое кол-во ошибок
		workerMaxCount = int32(n)   // максимальное возможное кол-во воркеров
		errCount       atomic.Int32 // текущее кол-во ошибок
		workerCount    atomic.Int32 // текущее кол-во работающих воркеров
	)

	if len(tasks) < n {
		workerMaxCount = int32(len(tasks))
	}

	for _, task := range tasks {
		// если количество работающих тасков меньше workerMaxCount -> запускаем таск
		// иначе, ждем 10 млсек
		for workerCount.Load() >= workerMaxCount {
			time.Sleep(time.Millisecond * 10)
		}

		if errCount.Load() >= errMaxCount { // если кол-во ошибок больше допустимого
			// прекращаем запускать воркеры
			break
		}

		workerCount.Add(1)
		wg.Add(1)
		go func(t Task) {
			defer func() {
				workerCount.Add(-1)
				wg.Done()
			}()

			err := t()
			if err != nil {
				errCount.Add(1)
			}
		}(task)
	}

	wg.Wait()

	if errCount.Load() >= errMaxCount {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// Run Tasks with channels.
func RunChan(tasks []Task, n, m int) error {
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

	tasksCh := make(chan Task)                  // канал задач.
	errorCh := make(chan error, workerMaxCount) // канал ошибок.
	wg := sync.WaitGroup{}

	defer func() {
		close(tasksCh)
		wg.Wait()
		close(errorCh)
	}()

	// запускаем воркеры
	for i := 0; i < workerMaxCount; i++ {
		wg.Add(1)
		go doWork(&wg, tasksCh, errorCh)
	}

	errCount := 0 // текущее кол-во ошибок

	for i := 0; i < len(tasks); {
		select {
		case <-errorCh:
			errCount++
			if errCount >= m {
				return ErrErrorsLimitExceeded
			}
		case tasksCh <- tasks[i]:
			i++
		}
	}

	return nil
}

// tasks - канал задач.
// errs - канал ошибок.
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
