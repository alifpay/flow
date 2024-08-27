package flow

import (
	"context"
	"fmt"
)

type Flow struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	StartNode *Node  `json:"start_node"`
}

// Node of flow
type Node struct {
	Rules     []Condition `json:"rules,omitempty"`
	Task      *Task       `json:"task,omitempty"`
	TrueNode  *Node       `json:"true_node,omitempty"`
	FalseNode *Node       `json:"false_node,omitempty"`
}

// Id of task function
type Task struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Process flow node
// if rules are correct go to true node
// if rules are not valid go to false node
func (n *Node) Process(ctx context.Context, input map[string]interface{}, errString string) error {
	// Если есть правила, проверяем их
	if len(n.Rules) != 0 {
		valid, errStr := Validate(n.Rules, input)
		if !valid {
			// Если есть узел False, выполняем его
			if n.FalseNode != nil {
				return n.FalseNode.Process(ctx, input, errStr)
			}
			return fmt.Errorf("validation failed: %s", errStr)
		}
	}
	// Если есть задача, выполняем ее
	if n.Task != nil {
		err := runTask(ctx, n.Task, input, errString)
		if err != nil {
			return err
		}
	}
	if n.TrueNode == nil {
		return nil
	}
	// Если есть следующий узел, выполняем его
	return n.TrueNode.Process(ctx, input, "")
}
