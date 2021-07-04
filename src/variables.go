package fzfyml

import (
	//"fmt"
	"path/filepath"
	"strings"
)

type Variables struct {
	vars map[string]string
}

func (variables *Variables) init(vars map[string]interface{}, ymlPath string) {
	variables.vars = map[string]string{}
	variables.vars["yml_dir"] = filepath.Dir(ymlPath)
	variables.initFromYml(vars)
}

func (variables *Variables) initFromYml(vars map[string]interface{}) {
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
