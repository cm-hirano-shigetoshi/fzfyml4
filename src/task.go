package fzfyml

import (
	"fmt"
	"os"
	"os/exec"
)

type Task struct {
	source         string
	variables      Variables
	options        Options
	postOperations PostOperations
}

func (task *Task) init(yml interface{}, ymlPath string) {
	baseTask := yml.(map[string]interface{})["base_task"].(map[string]interface{})
	task.source = baseTask["source"].(string)
	if _, ok := baseTask["variables"]; ok {
		task.variables.init(baseTask["variables"].(map[string]interface{}), ymlPath)
	}
	if _, ok := baseTask["options"]; ok {
		task.options.init(baseTask["options"].([]interface{}))
	}
	if _, ok := baseTask["post_operations"]; ok {
		task.postOperations.init(baseTask["post_operations"].([]interface{}))
	}
}

func (task *Task) run() Result {
	var result Result
	result.Init(task.execFzf(task.getExecuteCommand()))
	return result
}

func (task *Task) test(answer string) bool {
	response := task.getExecuteCommand()
	if answer == response {
		return true
	} else {
		fmt.Println(answer)
		fmt.Println(response)
		return false
	}
}

func (task *Task) isAppEnd(result Result) bool {
	return true
}

func (task *Task) execFzf(command string) string {
	cmd_exec := exec.Command("sh", "-c", command)
	cmd_exec.Stderr = os.Stderr
	out, _ := cmd_exec.Output()
	if len(out) > 0 {
		out = out[:len(out)-1]
		return string(out)
	} else {
		return ""
	}
}

func (task *Task) getExecuteCommand() string {
	optionText := task.options.getOptionText()
	optionText = task.variables.expand(optionText)
	return task.source + " | fzf " + optionText
}
