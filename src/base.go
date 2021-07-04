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
	task.init(yml)
	for {
		result := task.run()
		if task.isAppEnd(result) {
			fmt.Println(result.output)
			break
		}
	}
}

func Test(ymlPath string) {
	yml := getYml(ymlPath)
	tests := yml.(map[string]interface{})["test"].([]interface{})
	var task Task
	task.init(yml)

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
