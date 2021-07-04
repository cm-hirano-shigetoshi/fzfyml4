package fzfyml

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
)

func Run(ymlPath string) {
	yml := getYml(ymlPath)
	baseTask := yml.(map[string]interface{})["base_task"].(map[string]interface{})
	taskSwitch := map[string]interface{}{}
	if _, ok := yml.(map[string]interface{})["task_switch"]; ok {
		for key, val := range yml.(map[string]interface{})["task_switch"].(map[string]interface{}) {
			taskSwitch[key] = val
		}
	}
	var task Task
	task.init(baseTask, ymlPath)
	for {
		result := task.run()
		if newTask, ok := taskSwitch[result.key]; ok {
			task.update(newTask.(map[string]interface{}))
		} else {
			fmt.Print(task.postOperations.apply(result))
			break
		}
	}
}

func Test(ymlPath string) {
	yml := getYml(ymlPath)
	baseTask := yml.(map[string]interface{})["base_task"].(map[string]interface{})
	var task Task
	task.init(baseTask, ymlPath)

	tests := yml.(map[string]interface{})["test"].([]interface{})
	test := tests[0]
	if !task.test(test.(map[string]interface{})["answer"].(string)) {
		log.Fatal("test failed!")
	}
	for _, test := range tests[1:] {
		fmt.Println(task.test(test.(map[string]interface{})["answer"].(string)))
	}
}

func getYml(ymlPath string) interface{} {
	ymlData, _ := ioutil.ReadFile(ymlPath)
	var yml interface{}
	if err := yaml.Unmarshal([]byte(ymlData), &yml); err != nil {
		log.Fatal(err)
	}
	return yml
}
