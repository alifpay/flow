package flow

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EditCondition(db *pgxpool.Pool, jsonData []byte) error {
	var rules []Condition
	err := json.Unmarshal(jsonData, &rules)
	if err != nil {
		return err
	}
	for _, rule := range rules {
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
