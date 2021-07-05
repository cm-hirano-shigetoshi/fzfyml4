package fzfyml

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
)

func getYml(ymlPath string) interface{} {
	ymlData, _ := ioutil.ReadFile(ymlPath)
	var yml interface{}
	if err := yaml.Unmarshal([]byte(ymlData), &yml); err != nil {
		log.Fatal(err)
	}
	return yml
}

func initTask(yml interface{}, ymlPath string, args []string) (Task, map[string]interface{}) {
	baseTask := yml.(map[string]interface{})["base_task"].(map[string]interface{})
	taskSwitch := map[string]interface{}{}
	if _, ok := yml.(map[string]interface{})["task_switch"]; ok {
		for key, val := range yml.(map[string]interface{})["task_switch"].(map[string]interface{}) {
			taskSwitch[key] = val
		}
	}
	var task Task
	task.init(baseTask, ymlPath, args)
	return task, taskSwitch
}

func Run(ymlPath string, args []string) {
	yml := getYml(ymlPath)
	task, taskSwitch := initTask(yml, ymlPath, args)
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

func Test(ymlPath string, args []string) {
	yml := getYml(ymlPath)
	task, taskSwitch := initTask(yml, ymlPath, args)
	tests := yml.(map[string]interface{})["test"].([]interface{})
	test := tests[0]
	if !task.test(test.(map[string]interface{})["answer"].(string)) {
		log.Fatal("test failed!")
	}
	for _, test := range tests[1:] {
		var result Result
		result.initFromYml(test.(map[string]interface{})["result"].(map[string]interface{}))
		if newTask, ok := taskSwitch[result.key]; ok {
			task.update(newTask.(map[string]interface{}))
			if !task.test(test.(map[string]interface{})["answer"].(string)) {
				log.Fatal("test failed!")
			}
		} else {
			fmt.Print(task.postOperations.apply(result))
			break
		}
	}
}
