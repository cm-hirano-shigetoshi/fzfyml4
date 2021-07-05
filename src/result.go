package fzfyml

import (
	"strings"
)

type Result struct {
	query  string
	key    string
	output []string
}

func (result *Result) init(text string) {
	sp := strings.Split(text, "\n")
	result.query = sp[0]
	result.key = sp[1]
	result.output = sp[2:]
}

func (result *Result) initFromYml(obj map[string]interface{}) {
	result.query = obj["query"].(string)
	result.key = obj["key"].(string)
	result.output = []string{}
	for _, t := range obj["output"].([]interface{}) {
		result.output = append(result.output, t.(string))
	}
}
