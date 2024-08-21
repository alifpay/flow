package flow

import (
	"fmt"
	"reflect"
)

var dataInputs map[string]string

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
	Not      any      `json:"not,omitempty"`
}

type Condition struct {
	Type          ConditionType         `json:"type"`
	ErrMessage    string                `json:"errMessage"`
	Validation    map[string]Validation `json:"validation"`
	SubConditions *Condition            `json:"subConditions,omitempty"`
}

// data is map of inputs to validate
// rules is array of conditions to validate
func Validate(rules []Condition, data map[string]any) (valid bool, errStr string) {
	for _, rule := range rules {
		valid, errStr = conditionIsValid(rule, data)
		if !valid && rule.Type == AND {
			errStr = rule.ErrMessage + ": " + errStr
			return
		} else if valid && rule.Type == OR {
			return
		}
		// Если есть подусловия, проверяем их
		if rule.SubConditions != nil {
			valid, errStr = conditionIsValid(*rule.SubConditions, data)
			if !valid && rule.Type == AND {
				errStr = rule.ErrMessage + ": " + errStr
				return
			} else if valid && rule.Type == OR {
				return
			}
		}
		// rule type OR
		if !valid {
			errStr = rule.ErrMessage + ": " + errStr
			return
		}
	}
	return true, ""
}

func conditionIsValid(rule Condition, data map[string]any) (bool, string) {
	orValid := true
	for field, validation := range rule.Validation {
		value, exists := data[field]
		if !exists && validation.Required && rule.Type == AND {
			return false, fmt.Sprintf("%s обязателен", dataInputs[field])
		}
		if !exists {
			continue
		}

		// Проверка на минимальное значение
		if validation.Min != nil {
			if v, ok := toFloat(value); ok {
				if v < *validation.Min {
					if rule.Type == AND {
						return false, fmt.Sprintf("%s должен быть не менее %v, значение: %v", dataInputs[field], validation.Min, value)
					}
					orValid = false
				}
			}
		}
		// Проверка на максимальное значение
		if validation.Max != nil {
			if v, ok := toFloat(value); ok {
				if v > *validation.Max {
					if rule.Type == AND {
						return false, fmt.Sprintf("%s должен быть не более %v, значение: %v", dataInputs[field], validation.Max, value)
					}
					orValid = false
				}
			}
		}
		// Дополнительные проверки можно добавить здесь
		if validation.Equal != nil && !equal(value, validation.Equal) {
			if rule.Type == AND {
				return false, fmt.Sprintf("%s должен быть равен %v, значение: %v", dataInputs[field], validation.Equal, value)
			}
			orValid = false
		}

		if validation.Any != nil && !fnAny(validation.Any, value) {
			if rule.Type == AND {
				return false, fmt.Sprintf("%s должен быть в %v, значение: %v", field, validation.Any, value)
			}
			orValid = false
		}

		if validation.Not != nil && !fnNot(validation.Not, value) {
			if rule.Type == AND {
				return false, fmt.Sprintf("%s не должен быть в %v, значение: %v", field, validation.Not, value)
			}
			orValid = false
		}

		if rule.Type == OR && orValid {
			return true, ""
		}
	}
	return true, ""
}

// The ANY operator evaluates to true if equal for at least one value in the slice of valid.
// valid - slice of valid values
func fnAny(valid, value any) bool {
	v := reflect.ValueOf(valid)
	if v.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < v.Len(); i++ {
		element := v.Index(i).Interface()
		switch val := value.(type) {
		case int:
			if elementFloat, ok := element.(float64); ok && elementFloat == float64(val) {
				return true
			}
		case string:
			if elementString, ok := element.(string); ok && elementString == val {
				return true
			}
		default:
			if reflect.DeepEqual(element, value) {
				return true
			}
		}
	}
	return false
}

func fnNot(valid, value any) bool {
	v := reflect.ValueOf(valid)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return false
	}
	for i := 0; i < v.Len(); i++ {
		element := v.Index(i).Interface()
		switch val := value.(type) {
		case int:
			if elementFloat, ok := element.(float64); ok && elementFloat == float64(val) {
				return false
			}
		case string:
			if elementString, ok := element.(string); ok && elementString == val {
				return false
			}
		default:
			if reflect.DeepEqual(element, value) {
				return false
			}
		}
	}
	return true
}

func toFloat(val any) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func equal(a, b any) bool {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va == vb
		} else if vb, ok := b.(float64); ok {
			return float64(va) == vb
		}
	case float64:
		if vb, ok := b.(int); ok {
			return va == float64(vb)
		} else if vb, ok := b.(float64); ok {
			return va == vb
		}
	}
	return reflect.DeepEqual(a, b)
}
