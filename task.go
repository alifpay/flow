package flow

import (
	"context"
	"fmt"
)

type TaskFunc func(ctx context.Context, data map[string]any) error

var taskFuncs = map[string]TaskFunc{}

// Register registers a task function by name.
func Register(id string, fn TaskFunc) {
	taskFuncs[id] = fn
}

func runTask(ctx context.Context, t *Task, data map[string]any) error {
	if fn, found := taskFuncs[t.Id]; found {
		return fn(ctx, data)
	}
	return fmt.Errorf("функция %q не найдена", t.Name)
}
