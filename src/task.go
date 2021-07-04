package fzfyml

import (
	"os"
	"os/exec"
)

type Task struct {
	source string
}

func (task *Task) Init(yml interface{}) {
	baseTask := yml.(map[string]interface{})["base_task"].(map[string]interface{})
	task.source = baseTask["source"].(string)
}

func (task *Task) Run() Result {
	var result Result
	result.Init(task.execFzf(task.getExecuteCommand()))
	return result
}

func (task *Task) Test(answer string) bool {
	response := task.getExecuteCommand()
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
	return task.source + " | fzf --print-query --expect='ctrl-c'"
}
