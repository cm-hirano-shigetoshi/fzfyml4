package fzfyml

import (
	//"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	//"os"
)

type Variables struct {
	vars         map[string]string
	expandedVars map[string]string
}

func (variables *Variables) init(ymlPath string, vars map[string]interface{}) {
	variables.vars = map[string]string{}
	variables.expandedVars = map[string]string{}
	variables.vars["yml_dir"] = filepath.Dir(ymlPath)
	variables.updateFromYml(vars)
}

func (variables *Variables) update(vars map[string]interface{}) {
	variables.updateFromYml(vars)
}

func (variables *Variables) updateFromYml(vars map[string]interface{}) {
	for key, val := range vars {
		variables.vars[key] = val.(string)
	}
}

func (variables *Variables) expand(text string) string {
	variables.updateExpandedVars()
	for key, val := range variables.expandedVars {
		if strings.Contains(text, "{{"+key+"}}") {
			text = strings.Replace(text, "{{"+key+"}}", val, -1)
		}
	}
	return text
}

func (variables *Variables) updateExpandedVars() {
	expandedVars := map[string]string{}
	for key, val := range variables.vars {
		val = variables.expandShellCommand(val)
		expandedVars[key] = val
	}
	variables.expandedVars = expandedVars
}

func (variables *Variables) expandShellCommand(targetVal string) string {
	for key, _ := range variables.expandedVars {
		targetVal = strings.Replace(targetVal, "{{"+key+"}}", variables.expandedVars[key], -1)
	}
	targetVal = "echo \"" + targetVal + "\""
	cmd_exec := exec.Command("sh", "-c", targetVal)
	out, _ := cmd_exec.Output()
	if len(out) > 0 {
		out = out[:len(out)-1]
		return string(out)
	} else {
		return ""
	}
}
