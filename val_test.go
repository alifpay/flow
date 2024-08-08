package flow

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCondition(t *testing.T) {
	var rule []Condition
	jsonData := `
    [{
        "type": "AND",
        "errMessage": "Main condition error",
        "validation": {
            "age": {
                "required": true,
                "max": 30.0
            }
        },
        "subConditions": {
            "type": "OR",
            "errMessage": "Sub condition error",
            "validation": {
                "amount": {
                    "required": true,
                    "min": 10.0,
                    "max": 3000.0
                }
            }
        }
    }]`

	err := json.Unmarshal([]byte(jsonData), &rule)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Вывод результата десериализации
	fmt.Printf("Condition: %+v\n", rule)
	fmt.Printf("Condition: %+v\n", rule[0].SubConditions)
}

func TestValidate(t *testing.T) {

}