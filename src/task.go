package fzfyml

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Task struct {
	source         string
	variables      Variables
	binds          Binds
	preview        Preview
	options        Options
	postOperations PostOperations
	switchExpects  []string
}

func (task *Task) init(baseTask map[string]interface{}, ymlPath string, switchExpects []string, args []string) {
	task.source = baseTask["source"].(string)
	variables, _ := baseTask["variables"].(map[string]interface{})
	task.variables.init(ymlPath, args, variables)
	if _, ok := baseTask["binds"]; ok {
		task.binds.init(baseTask["binds"].(map[string]interface{}))
	}
	if _, ok := baseTask["preview"]; ok {
		task.preview.init(baseTask["preview"].(map[string]interface{}))
	}
	if _, ok := baseTask["options"]; ok {
		task.options.init(baseTask["options"].([]interface{}))
	}
	if _, ok := baseTask["post_operations"]; ok {
		task.postOperations.init(baseTask["post_operations"].(map[string]interface{}))
	}
	task.switchExpects = switchExpects
}

func (task *Task) update(newTask map[string]interface{}) {
	if _, ok := newTask["source"]; ok {
		task.source = newTask["source"].(string)
	}
	if _, ok := newTask["variables"]; ok {
		task.variables.update(newTask["variables"].(map[string]interface{}))
	}
	if _, ok := newTask["binds"]; ok {
		task.binds.update(newTask["binds"].(map[string]interface{}))
	}
}

func (task *Task) run() Result {
	var result Result
	result.init(task.execFzf(task.getExecuteCommand("run")))
	return result
}

func (task *Task) test(answer string) bool {
	response := task.getExecuteCommand("test")
	if answer == response {
		return true
	} else {
		fmt.Println(answer)
		fmt.Println(response)
		return false
	}
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

func (task *Task) getExecuteCommand(mode string) string {
	bindList := task.binds.getBindList()
	preview := task.preview.getPreviewText()
	optionList := task.options.getOptionList()
	expectList := task.getExpectList()
	mondatoryList := []string{"--print-query"}
	if mode == "test" {
		sort.Strings(bindList)
		sort.Strings(optionList)
		sort.Strings(expectList)
		sort.Strings(mondatoryList)
	}
	command := task.source + " | fzf " + strings.Join(bindList, " ") + " " + preview + " " + strings.Join(optionList, " ") + " --expect=" + strings.Join(expectList, ",") + " " + strings.Join(mondatoryList, " ")
	command = task.variables.expand(command)
	//fmt.Println(command+"\n")
	return command
}

func (task *Task) getExpectList() []string {
	expects := []string{}
	for _, key := range task.postOperations.getExpects() {
		expects = append(expects, key)
	}
	for _, key := range task.switchExpects {
		expects = append(expects, key)
	}
	expects = append(expects, strings.Split("enter,esc,ctrl-c,ctrl-d,ctrl-g,ctrl-q,ctrl-z", ",")...)
	expects = uniqueStringSlice(expects)
	return expects
}
