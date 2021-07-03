package fzfyml

import (
	//"fmt"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
)

func Run(ymlPath string) {
	ymlData, _ := ioutil.ReadFile(ymlPath)
	var yml interface{}
	if err := yaml.Unmarshal([]byte(ymlData), &yml); err != nil {
		log.Fatal(err)
	}

	var task Task
	task.Init(yml)
	task.Run()
}
