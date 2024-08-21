package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EditFlow(ctx context.Context, db *pgxpool.Pool, jsonData []byte) error {
	var f Flow
	err := json.Unmarshal(jsonData, &f)
	if err != nil {
		return err
	}
	err = isNodeValid(f.StartNode)
	if err != nil {
		return err
	}
	cmd, err := db.Exec(ctx, `INSERT INTO flows (id, name, node) VALUES ($1, $2, $3) 
				  ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, node = EXCLUDED.node;`, f.Id, f.Name, f.StartNode)
	if err != nil {
		return err
	}
	if ra := cmd.RowsAffected(); ra != 1 {
		return errors.New("db rows affected is wrong")
	}
	return nil
}

func isNodeValid(node *Node) error {
	if len(node.Rules) > 0 {
		for _, rule := range node.Rules {
			err := isValidEditCondition(rule)
			if err != nil {
				return err
			}
			if rule.SubConditions != nil {
				err := isValidEditCondition(*rule.SubConditions)
				if err != nil {
					return err
				}
			}
		}
	}

	if node.Task != nil {
		if _, found := taskFuncs[node.Task.Id]; !found {
			return fmt.Errorf("task function %q not found", node.Task.Name)
		}
	}
	if node.TrueNode != nil {
		err := isNodeValid(node.TrueNode)
		if err != nil {
			return err
		}
	}
	if node.FalseNode != nil {
		err := isNodeValid(node.FalseNode)
		if err != nil {
			return err
		}
	}
	return nil
}

// validate input condition
func isValidEditCondition(rule Condition) error {
	if rule.Type != AND && rule.Type != OR {
		return fmt.Errorf("unknown condition type: %s", rule.Type)
	}
	rule.ErrMessage = strings.TrimSpace(rule.ErrMessage)
	if len(rule.ErrMessage) < 5 {
		return fmt.Errorf("errMessage is too short: %s", rule.ErrMessage)
	}
	for field, validation := range rule.Validation {
		if validation.Min != nil || validation.Max != nil {
			if validation.Equal != nil {
				return fmt.Errorf("field %s: can't have equal and min/max at the same time", field)
			}
			if validation.Any != nil {
				return fmt.Errorf("field %s: can't have any and min/max at the same time", field)
			}
			if validation.Not != nil {
				return fmt.Errorf("field %s: can't have not and min/max at the same time", field)
			}
		} else if validation.Equal != nil {
			if validation.Any != nil {
				return fmt.Errorf("field %s: can't have any and equal at the same time", field)
			}
			if validation.Not != nil {
				return fmt.Errorf("field %s: can't have not and equal at the same time", field)
			}
		} else if validation.Any != nil && validation.Not != nil {
			return fmt.Errorf("field %s: can't have not and any at the same time", field)
		}
	}
	return nil
}

func GetFlow(ctx context.Context, db *pgxpool.Pool, id string) (*Flow, error) {
	var f Flow
	err := db.QueryRow(ctx, `SELECT name, node FROM flows WHERE id = $1`, id).Scan(&f.Name, &f.StartNode)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func GetDataInputs(ctx context.Context, db *pgxpool.Pool) (map[string]string, error) {
	dataInputs := make(map[string]string)
	rows, err := db.Query(ctx, "SELECT id, name FROM data_inputs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		dataInputs[id] = name
	}
	return dataInputs, nil
}

func EditDataInput(ctx context.Context, db *pgxpool.Pool, jsonData []byte) error {
	var dataInput struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	err := json.Unmarshal(jsonData, &dataInput)
	if err != nil {
		return err
	}
	dataInput.Id = strings.TrimSpace(dataInput.Id)
	dataInput.Name = strings.TrimSpace(dataInput.Name)
	if len(dataInput.Id) < 2 {
		return errors.New("id is not valid")
	}
	if len(dataInput.Name) < 2 {
		return errors.New("name is not valid")
	}
	cmd, err := db.Exec(ctx, `INSERT INTO data_inputs (id, name) VALUES ($1, $2) 
				  ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;`, dataInput.Id, dataInput.Name)
	if err != nil {
		return err
	}
	if ra := cmd.RowsAffected(); ra != 1 {
		return errors.New("db rows affected is wrong")
	}
	return nil
}

func EditFunction(ctx context.Context, db *pgxpool.Pool, jsonData []byte) error {
	var fn struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	err := json.Unmarshal(jsonData, &fn)
	if err != nil {
		return err
	}
	fn.Id = strings.TrimSpace(fn.Id)
	fn.Name = strings.TrimSpace(fn.Name)
	if len(fn.Id) < 2 {
		return errors.New("id is not valid")
	}
	if len(fn.Name) < 2 {
		return errors.New("name is not valid")
	}
	cmd, err := db.Exec(ctx, `INSERT INTO functions (id, name) VALUES ($1, $2) 
				  ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;`, fn.Id, fn.Name)
	if err != nil {
		return err
	}
	if ra := cmd.RowsAffected(); ra != 1 {
		return errors.New("db rows affected is wrong")
	}
	return nil
}

func GetFunctions(ctx context.Context, db *pgxpool.Pool) (map[string]string, error) {
	fns := make(map[string]string)
	rows, err := db.Query(ctx, "SELECT id, name FROM functions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		fns[id] = name
	}
	return fns, nil
}
