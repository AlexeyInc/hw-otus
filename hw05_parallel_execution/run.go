package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var result = make(chan error)

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	goroutineCounter := make(chan struct{}, n)

	doneCh := make(chan error)
	errCh := make(chan error)
	quitCh := make(chan bool)
	tasksCount := len(tasks)

	go checkTasksCompletion(doneCh, errCh, quitCh, tasksCount, m)

	for _, curTask := range tasks {
		select {
		case goroutineCounter <- struct{}{}:
			go executeTask(curTask, doneCh, errCh, quitCh, goroutineCounter)
		case <-quitCh:
			break
		}
	}
	return <-result
}

func checkTasksCompletion(doneCh chan error, errCh chan error, quitCh chan bool, taskCount int, maxErrors int) {
	tasksCount := 0
	errCount := 0
	for {
		select {
		case <-errCh:
			errCount++
			if maxErrors > 0 && errCount >= maxErrors {
				finish(quitCh, ErrErrorsLimitExceeded)
				return
			}

			tasksCount++
			if tasksCount == taskCount {
				finish(quitCh, nil)
				return
			}
		case <-doneCh:
			tasksCount++
			if tasksCount == taskCount {
				finish(quitCh, nil)
				return
			}
		}
	}
}

func finish(quit chan bool, err error) {
	close(quit)
	result <- err
}

func executeTask(taskFunc Task,
	doneCh chan error,
	errCh chan error,
	quitCh chan bool,
	goroutineCounter chan struct{}) {
	err := taskFunc() //nolint

	ch := doneCh
	if err != nil {
		ch = errCh
	}

	select {
	case ch <- err:
	case <-quitCh:
	}

	<-goroutineCounter
}
