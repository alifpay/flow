# Flow

### Use Case: Dynamic Stop List

Этот пример использования включает динамический стоп-лист, который оценивает серию условий для определения потока выполнения. Стоп-лист содержит различные условия, проверяющие входные значения. На основе этих условий система определяет поток и выполняет соответствующую функцию, в конечном итоге возвращая результат.

### Описание

Реализация системы задач или условий в Golang, которая может быть настроена и сохранена в базе данных, аналогично блок-схеме.

### Основная структура

1. Определите структуры для задач и условий.
2. Реализуйте способ соединения задач и условий.
3. Создайте методы для оценки условий и выполнения задач.
4. Разработайте систему для сохранения и загрузки потока из базы данных.

### Поддерживаемые условия

- Логические операции: AND, OR
- Операции сравнения: больше, меньше


AND[rule1, rule2, OR[rule3, rule4]]

OR[rule1, rule2, AND[rule3, rule4]]

rule[age] = struct{
    requered,
    min,
    max,
    equal,
    any[],
    not[]
}


example 

```json
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
    }]
```

Flow 

```Go
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
```

Todo :  
database +
example  
test  

```Go

func PrintHello(ctx context.Context, data map[string]any, errStr string) error {
    // Check if the context has been canceled
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Continue processing
    }
    fmt.Println("Hello world!", errStr)
    // Example of a successful operation
    return nil
}

    //init datainputs
    err := flow.InitDataInputs(ctx, db)
    if err != nil{
        log.Fatal(err)
    }

    // register your task function
    flow.RegisterTask("PrintHello", PrintHello)
```

