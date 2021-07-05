package fzfyml

import (
	"strings"
)

type PostOperations struct {
	list map[string]interface{}
}

func (operations *PostOperations) init(operationList map[string]interface{}) {
	operations.list = operationList
}

func (operations *PostOperations) apply(result Result) string {
	output := strings.Join(result.output, "\n")
	if ope, ok := operations.list[result.key].([]interface{}); ok {
		for _, operation := range ope {
			switch operation.(type) {
			case string:
				if operation.(string) == "join" {
					output = strings.Replace(output, "\n", " ", -1)
				}
			}
		}
	}
	return output
}
