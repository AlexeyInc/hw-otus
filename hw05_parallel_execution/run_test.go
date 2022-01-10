package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t,
			atomic.LoadInt32(&runTasksCount), int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks with errors limit in last task", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		workersCount := 5
		maxErrorsCount := 1

		var runTasksCount int32
		taskSleep := time.Millisecond * 10
		tasks = InitTasks(tasks, tasksCount-maxErrorsCount, taskSleep, &runTasksCount, nil)

		err := fmt.Errorf("error from task")
		taskSleep = time.Millisecond * 100
		tasks = InitTasks(tasks, maxErrorsCount, taskSleep, &runTasksCount, err)

		result := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(result, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("tasks with error in last task, but not error limit", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		workersCount := 5
		maxErrorsCount := 2

		taskSleep := time.Millisecond * 50
		var runTasksCount int32
		tasks = InitTasks(tasks, tasksCount-1, taskSleep, &runTasksCount, nil)

		err := fmt.Errorf("error from task")
		taskSleep = time.Millisecond * 500
		tasks = InitTasks(tasks, maxErrorsCount-1, taskSleep, &runTasksCount, err)

		result := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, !errors.Is(result, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("amount of tasks less than workers", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		workersCount := 15
		maxErrorsCount := 2

		taskSleep := time.Millisecond * 100
		var runTasksCount int32
		tasks = InitTasks(tasks, tasksCount, taskSleep, &runTasksCount, nil)

		result := Run(tasks, workersCount, maxErrorsCount)

		require.True(t, !errors.Is(result, ErrErrorsLimitExceeded))
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("check count of running goroutines", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 10
		maxErrorsCount := 2

		taskSleep := time.Millisecond * 100
		var runTasksCount int32
		tasks = InitTasks(tasks, tasksCount, taskSleep, &runTasksCount, nil)

		Run(tasks, workersCount, maxErrorsCount)
		numGoroutine := runtime.NumGoroutine()

		require.LessOrEqual(t, numGoroutine, workersCount+1, "extra goroutine were started")
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("if m <= 0 ignore errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 10
		maxErrorsCount := 0

		var runTasksCount int32
		halfOfTasks := tasksCount / 2
		taskSleep := time.Microsecond * 10
		tasks = InitTasks(tasks, halfOfTasks, taskSleep, &runTasksCount, nil)

		err := fmt.Errorf("error from task")
		tasks = InitTasks(tasks, halfOfTasks, taskSleep, &runTasksCount, err)

		result := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, result)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}

func InitTasks(tasks []Task, count int, taskSleep time.Duration, runTasksCount *int32, err error) []Task {
	for i := 0; i < count; i++ {
		tasks = append(tasks, func() error {
			time.Sleep(taskSleep)
			atomic.AddInt32(runTasksCount, 1)
			return err
		})
	}
	return tasks
}
