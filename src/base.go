package fzfyml

import (
	"bufio"
	"fmt"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
	"os"
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
	switchExpects := []string{}
	taskSwitch := map[string]interface{}{}
	if _, ok := yml.(map[string]interface{})["task_switch"]; ok {
		for key, val := range yml.(map[string]interface{})["task_switch"].(map[string]interface{}) {
			switchExpects = append(switchExpects, key)
			taskSwitch[key] = val
		}
	}
	var task Task
	task.init(baseTask, ymlPath, switchExpects, args)
	return task, taskSwitch
}

func execute(ymlPath string, args []string) (int, error) {
	yml := getYml(ymlPath)
	task, taskSwitch := initTask(yml, ymlPath, args)
	var result Result
	result.init(task.run(nil))
	for {
		if newTask, ok := taskSwitch[result.key]; ok {
			// task_switch後は元タスクのquery指定を無効にする
			task.variables.updateResult(result)
			task.options.deleteQuery()
			task.update(newTask.(map[string]interface{}))
		} else {
			fmt.Print(task.postOperations.apply(result))
			break
		}
		result.init(task.run(result.query))
	}
	return 0, nil
}

func test(ymlPath string) (int, error) {
	yml := getYml(ymlPath)
	tests := yml.(map[string]interface{})["test"].([]interface{})
	test := tests[0]
	args := []string{}
	if argsList, ok := test.(map[string]interface{})["args"].([]interface{}); ok {
		for _, arg := range argsList {
			args = append(args, arg.(string))
		}
	}
	task, taskSwitch := initTask(yml, ymlPath, args)
	if !task.test(nil, test.(map[string]interface{})["answer"].(string)) {
		return 1, fmt.Errorf("test failed!")
	}
	for _, test := range tests[1:] {
		var result Result
		result.initFromYml(test.(map[string]interface{})["result"].(map[string]interface{}))
		if newTask, ok := taskSwitch[result.key]; ok {
			// task_switch後は元タスクのquery指定を無効にする
			task.options.deleteQuery()
			task.update(newTask.(map[string]interface{}))
			if !task.test(result.query, test.(map[string]interface{})["answer"].(string)) {
				log.Fatal("test failed!")
			}
		} else {
			fmt.Print(task.postOperations.apply(result))
			break
		}
	}
	return 0, nil
}

func innerUntransformedOutput(textFilePath string, indexFilePath string) (int, error) {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	fmt.Println(sc.Text())
	sc.Scan()
	fmt.Println(sc.Text())
	lines := selectLine(textFilePath, indexFilePath)
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0, nil
}
