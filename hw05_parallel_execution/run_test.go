package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		var runTasksCount int32
		tasksCount := 50

		tasks, _ := createTasks(t, tasksCount, &runTasksCount, true)

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		var runTasksCount int32
		tasksCount := 50

		tasks, sumTime := createTasks(t, tasksCount, &runTasksCount, false)

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("negative or zero parameters", func(t *testing.T) {
		var err error

		err = Run([]Task{}, 10, 0)
		require.True(t, errors.Is(err, ErrErrorsLimitExceeded))

		err = Run([]Task{}, 10, -10)
		require.True(t, errors.Is(err, ErrErrorsLimitExceeded))

		err = Run([]Task{}, 0, 10)
		require.True(t, errors.Is(err, ErrCannotRunWorkers))

		err = Run([]Task{}, -10, 10)
		require.True(t, errors.Is(err, ErrCannotRunWorkers))
	})

	t.Run("count tasks less than count workers", func(t *testing.T) {
		var runTasksCount int32
		tasksCount := 10

		tasks, sumTime := createTasks(t, tasksCount, &runTasksCount, false)

		workersCount := 25
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks with and without errors, result - error", func(t *testing.T) {
		var runTasksCount int32

		tasks, _ := createMixedTasks(t, 10, 10, &runTasksCount)

		workersCount := 3
		maxErrorsCount := 15
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount+10), "extra tasks were started")
	})

	t.Run("tasks with and without errors, result - success", func(t *testing.T) {
		var runTasksCount int32

		tasks, sumTime := createMixedTasks(t, 10, 50, &runTasksCount)

		workersCount := 3
		maxErrorsCount := 25

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(2*10+50), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func createTasks(t *testing.T, tasksCount int, runTasksCount *int32, returnErr bool) ([]Task, time.Duration) {
	t.Helper()
	tasks := make([]Task, 0, tasksCount)

	var sumTime time.Duration

	for i := 0; i < tasksCount; i++ {
		taskSleep := time.Millisecond * (time.Duration(rand.Intn(100)) + 1)
		sumTime += taskSleep
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			require.Eventually(t, func() bool { return true }, 2*taskSleep, taskSleep)
			atomic.AddInt32(runTasksCount, 1)
			if returnErr {
				return err
			}
			return nil
		})
	}

	return tasks, sumTime
}

func createMixedTasks(t *testing.T, withErrorCnt, withoutErrorCnt int, runTasksCount *int32) ([]Task, time.Duration) {
	t.Helper()

	var tasks []Task
	var sumTime time.Duration

	tasksTmp, sumTimeTmp := createTasks(t, withErrorCnt, runTasksCount, true)
	tasks = append(tasks, tasksTmp...)
	sumTime += sumTimeTmp

	tasksTmp, sumTimeTmp = createTasks(t, withoutErrorCnt, runTasksCount, false)
	tasks = append(tasks, tasksTmp...)
	sumTime += sumTimeTmp

	tasksTmp, sumTimeTmp = createTasks(t, withErrorCnt, runTasksCount, true)
	tasks = append(tasks, tasksTmp...)
	sumTime += sumTimeTmp

	return tasks, sumTime
}
