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

// can you write test for Validate function?
func TestValidate(t *testing.T) {
	//parse json for conditions
	var rule []Condition
	//give example of json data
	jsonData := `
    [{
        "type": "AND",
        "errMessage": "Main condition error",
        "validation": {
            "age": {
                "required": true,
                "max": 30.0
            },
            "test": {
                "not": [3,4]
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
	fmt.Printf("Condition: %+v\n", rule)
	data := map[string]any{
		"age":    25,
		"amount": 100,
		"test":   7,
	}

	valid, errStr := Validate(rule, data)
	if !valid {
		t.Fatalf("Validation failed: %v", errStr)
	}
}

type Test struct {
	Any   any `json:"any,omitempty"`
	Not   any `json:"not,omitempty"`
	Equal any `json:"equal,omitempty"`
}

func TestNotJson(t *testing.T) {
	// Пример JSON данных
	jsonData := `{"not": [3, 4]}`

	// Десериализация JSON в структуру
	var test Test
	err := json.Unmarshal([]byte(jsonData), &test)
	if err != nil {
		fmt.Println("Ошибка при десериализации:", err)
		return
	}

	fmt.Printf("Десериализованная структура: %+v\n", test.Not)

	// Сериализация структуры обратно в JSON
	serializedData, err := json.Marshal(test)
	if err != nil {
		fmt.Println("Ошибка при сериализации:", err)
		return
	}

	fmt.Printf("Сериализованные данные: %s\n", string(serializedData))
}

// write test fnNot
func TestFnNot(t *testing.T) {
	// Пример JSON данных
	jsonData := `{"not": [3, 4]}`

	// Десериализация JSON в структуру
	var test Test
	err := json.Unmarshal([]byte(jsonData), &test)
	if err != nil {
		fmt.Println("Ошибка при десериализации:", err)
		return
	}

	fmt.Printf("Десериализованная структура: %+v\n", test.Not)

	// Проверка на соответствие условию
	if fnNot(test.Not, 3) {
		t.Fatalf("Ошибка: 3 в списке")
	}
	if !fnNot(test.Not, 5) {
		t.Fatalf("Ошибка: 5 нет в списке")
	}
}

func TestFnNotStr(t *testing.T) {
	// Пример JSON данных
	jsonData := `{"not": ["a", "b"]}`

	// Десериализация JSON в структуру
	var test Test
	err := json.Unmarshal([]byte(jsonData), &test)
	if err != nil {
		fmt.Println("Ошибка при десериализации:", err)
		return
	}

	fmt.Printf("Десериализованная структура: %+v\n", test.Not)

	// Проверка на соответствие условию
	if fnNot(test.Not, "a") {
		t.Fatalf("Ошибка: a в списке")
	}

	if !fnNot(test.Not, "c") {
		t.Fatalf("Ошибка: c нет в списке")
	}
}

func TestFnAny(t *testing.T) {
	// Пример JSON данных
	jsonData := `{"any": [3, 4,7]}`

	// Десериализация JSON в структуру
	var test Test
	err := json.Unmarshal([]byte(jsonData), &test)
	if err != nil {
		fmt.Println("Ошибка при десериализации:", err)
		return
	}

	fmt.Printf("Десериализованная структура: %+v\n", test.Any)

	// Проверка на соответствие условию
	if !fnAny(test.Any, 3) {
		t.Fatalf("Ошибка: 3 в списке")
	}

	if fnAny(test.Any, 5) {
		t.Fatalf("Ошибка: 5 нет в списке")
	}
}

func TestEqual(t *testing.T) {
	// Пример JSON данных
	jsonData := `{"equal": "dd"}`

	// Десериализация JSON в структуру
	var test Test
	err := json.Unmarshal([]byte(jsonData), &test)
	if err != nil {
		fmt.Println("Ошибка при десериализации:", err)
		return
	}

	fmt.Printf("Десериализованная структура: %T\n", test.Equal)

	// Проверка на соответствие условию
	if equal("ww", test.Equal) {
		t.Fatalf("Ошибка: 3 в списке")
	}

	if !equal("dd", test.Equal) {
		t.Fatalf("Ошибка: 5 нет в списке")
	}
}
