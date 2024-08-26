package flow

import (
	"context"
	"fmt"
)

type TaskFunc func(ctx context.Context, data map[string]any, errStr string) error

var taskFuncs = map[string]TaskFunc{}

// Register registers a task function by name.
func RegisterTask(id string, fn TaskFunc) {
	taskFuncs[id] = fn
}

func runTask(ctx context.Context, t *Task, data map[string]any, errStr string) error {
	if fn, found := taskFuncs[t.Id]; found {
		return fn(ctx, data, errStr)
	}
	return fmt.Errorf("функция %q не найдена", t.Name)
}
