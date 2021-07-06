package fzfyml

import (
	"strings"
	//"fmt"
)

type PostOperations struct {
	list map[string]interface{}
}

func (operations *PostOperations) init(operationList map[string]interface{}) {
	operations.list = operationList
}

func (operations *PostOperations) apply(result Result) string {
	text := strings.Join(result.output, "\n")
	if ope, ok := operations.list[result.key].([]interface{}); ok {
		for _, operation := range ope {
			switch operation.(type) {
			case string:
				if operation.(string) == "join" {
					text = strings.Replace(text, "\n", " ", -1)
				}
			case map[string]interface{}:
				if pipe, ok := operation.(map[string]interface{})["pipe"].(string); ok {
					text = execCommandWithInput(pipe, text)
				}
			}
		}
	}
	return text
}

func (operations *PostOperations) getExpects() []string {
	expects := []string{}
	for key, _ := range operations.list {
		expects = append(expects, key)
	}
	return expects

}
