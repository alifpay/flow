package flow

import (
	"fmt"
	"reflect"
)

// ConditionType определяет тип условия (AND или OR)
type ConditionType string

const (
	AND ConditionType = "AND"
	OR  ConditionType = "OR"
)

type Validation struct {
	Required bool     `json:"required,omitempty"`
	Min      *float64 `json:"min,omitempty"`
	Max      *float64 `json:"max,omitempty"`
	Equal    any      `json:"equal,omitempty"`
	Any      any      `json:"any,omitempty"`
	All      any      `json:"all,omitempty"`
}

type Condition struct {
	Type          ConditionType         `json:"type"`
	ErrMessage    string                `json:"errMessage"`
	Validation    map[string]Validation `json:"validation"`
	SubConditions *Condition            `json:"subConditions,omitempty"`
}

func Validate(rules []Condition, data map[string]any) (valid bool, errStr string) {
	for _, rule := range rules {
		valid, errStr = conditionIsValid(rule, data)
		if !valid && rule.Type == AND {
			return
		} else if valid && rule.Type == OR {
			return
		}
		// Если есть подусловия, проверяем их
		if rule.SubConditions != nil {
			subValid, subErr := conditionIsValid(*rule.SubConditions, data)
			if !subValid && rule.Type == AND {
				valid = false
				errStr = subErr
				return 
			}else if subValid && rule.Type == OR {
				return true, ""
			}
		}
	}

	return true, ""
}

func conditionIsValid(rule Condition, data map[string]any) (bool, string) {
	for field, validation := range rule.Validation {
		value, exists := data[field]
		if !exists && validation.Required {
			return false, fmt.Sprintf("%s is required", field)
		}
		if !exists {
			return false, ""
		}
		// Проверка на минимальное значение
		if validation.Min != nil {
			if v, ok := value.(float64); ok {
				if v < *validation.Min {
					return false, fmt.Sprintf("%s must be at least %v", field, validation.Min)
				}
			}
		}
		// Проверка на максимальное значение
		if validation.Max != nil {
			if v, ok := value.(float64); ok {
				if v > *validation.Max {
					return false, fmt.Sprintf("%s must be at most %v", field, validation.Max)
				}
			}
		}
		// Дополнительные проверки можно добавить здесь
		if validation.Equal != nil {
			if !reflect.DeepEqual(value, validation.Equal) {
				return false, fmt.Sprintf("%s must be equal to %v", field, validation.Equal)
			}
		}
		if validation.Any != nil {
			if !fnAny(validation.Any, value) {
				return false, fmt.Sprintf("%s must be in %v", field, validation.Any)
			}
		}

		if validation.All != nil {
			if !fnAll(validation.All, value) {
				return false, fmt.Sprintf("%s must be in %v", field, validation.All)
			}
		}
	}
	return true, ""
}

// valid is slice
func fnAny(valid, value any) bool {
	v := reflect.ValueOf(valid)
	if v.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(v.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

func fnAll(valid, value any) bool {
	v := reflect.ValueOf(valid)
	if v.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < v.Len(); i++ {
		if !reflect.DeepEqual(v.Index(i).Interface(), value) {
			return false
		}
	}
	return true
}
