package fzfyml

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	//"os"
)

type Variables struct {
	vars         map[string]string
	expandedVars map[string]string
}

func (variables *Variables) init(ymlPath string, args []string, vars map[string]interface{}) {
	variables.vars = map[string]string{}
	variables.expandedVars = map[string]string{}
	variables.vars["yml_dir"] = filepath.Dir(ymlPath)
	for i, arg := range args {
		variables.vars["arg"+fmt.Sprint(i+1)] = arg
	}
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

func (variables *Variables) updateResult(result Result) {
	variables.vars["query"] = result.query
	variables.vars["key"] = result.key
}

func (variables *Variables) expand(text string) string {
	variables.updateExpandedVars()

	keys := []string{}
	for k, _ := range variables.expandedVars {
		keys = append(keys, k)
	}
	parallel := strings.Join(keys, "|")
	pattern := regexp.MustCompile(`{{(` + parallel + `)}}`)
	matches := pattern.FindAllStringSubmatchIndex(text, -1)
	for len(matches) > 0 {
		for i := len(matches) - 1; i >= 0; i-- {
			start := matches[i][0]
			end := matches[i][1]
			key := text[start+2 : end-2]
			text = text[:start] + variables.expandedVars[key] + text[end:]
		}
		matches = pattern.FindAllStringSubmatchIndex(text, -1)
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
