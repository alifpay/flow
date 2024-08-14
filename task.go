package flow

import (
	"context"
	"fmt"
)

type TaskFunc func(ctx context.Context, input map[string]any) error

var taskFuncs = map[string]TaskFunc{}

// Register registers a task function by name.
func Register(id string, fn TaskFunc) {
	taskFuncs[id] = fn
}

func runTask(ctx context.Context, t *Task, input map[string]any) error {
	if fn, found := taskFuncs[t.Id]; found {
		return fn(ctx, input)
	}
	return fmt.Errorf("Task %q not found", t.Name)
}
