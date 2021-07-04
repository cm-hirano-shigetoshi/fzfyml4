package fzfyml

import (
	"fmt"
	"os"
	"os/exec"
)

type Task struct {
	source  string
	options Options
}

func (task *Task) Init(yml interface{}) {
	baseTask := yml.(map[string]interface{})["base_task"].(map[string]interface{})
	task.source = baseTask["source"].(string)
	task.options.InitFromYml(baseTask["options"].([]interface{}))
}

func (task *Task) Run() Result {
	var result Result
	result.Init(task.execFzf(task.getExecuteCommand()))
	return result
}

func (task *Task) Test(answer string) bool {
	response := task.getExecuteCommand()
	fmt.Println(response)
	return answer == response
}

func (task *Task) IsAppEnd(result Result) bool {
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
	optionText := task.options.GetOptionText()
	return task.source + " | fzf " + optionText
}
