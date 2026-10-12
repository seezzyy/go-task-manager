package main

import (
	"encoding/json"
	"fmt"
	"os"
	"task-manager/internal/model"
)

func main() {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var tasks []model.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))

}
