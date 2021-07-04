package fzfyml

import (
	"strings"
)

type PostOperations struct {
	list []interface{}
}

func (operations *PostOperations) init(operationList []interface{}) {
	operations.list = operationList
}

func (operations *PostOperations) apply(result Result) string {
	output := strings.Join(result.output, "\n")
	for _, operation := range operations.list {
		switch operation.(type) {
		case string:
			if operation.(string) == "join" {
				output = strings.Replace(output, "\n", " ", -1)
			}
		}
	}
	return output
}
