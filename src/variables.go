package fzfyml

import (
	//"fmt"
	"strings"
)

type Variables struct {
	vars map[string]string
}

func (variables *Variables) init(vars map[string]interface{}) {
	variables.vars = map[string]string{}
	for key, val := range vars {
		variables.vars[key] = val.(string)
	}
}

func (variables *Variables) expand(text string) string {
	for key, val := range variables.vars {
		if strings.Contains(text, "{{"+key+"}}") {
			text = strings.Replace(text, "{{"+key+"}}", val, -1)
		}
	}
	return text
}
