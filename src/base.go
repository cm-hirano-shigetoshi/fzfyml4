package fzfyml

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
)

func Run(ymlPath string) {
	yml := getYml(ymlPath)
	var task Task
	task.Init(yml)
	for {
		result := task.Run()
		if true {
			fmt.Println(result.output)
			break
		}
	}
}

func Test(ymlPath string) {
	yml := getYml(ymlPath)
	tests := yml.(map[string]interface{})["test"].([]interface{})
	var task Task
	task.Init(yml)
	fmt.Println(task.Test(tests[0].(map[string]interface{})["answer"].(string)))
}

func getYml(ymlPath string) interface{} {
	ymlData, _ := ioutil.ReadFile(ymlPath)
	var yml interface{}
	if err := yaml.Unmarshal([]byte(ymlData), &yml); err != nil {
		log.Fatal(err)
	}
	return yml
}
