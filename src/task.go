package fzfyml

import (
	"fmt"
	"os"
	"os/exec"
)

type Task struct {
	source string
}

func (t *Task) Init(yml interface{}) {
	baseTask := yml.(map[string]interface{})["base_task"].(map[string]interface{})
	t.source = baseTask["source"].(string)
}

func (t *Task) Run() {
	fmt.Print(fzf(t.source))
}

func fzf(source string) string {
	cmd := source + " | fzf"
	cmd_exec := exec.Command("sh", "-c", cmd)
	cmd_exec.Stderr = os.Stderr
	out, _ := cmd_exec.Output()
	if len(out) > 0 {
		out = out[:len(out)-1]
		return string(out)
	} else {
		return ""
	}
}
