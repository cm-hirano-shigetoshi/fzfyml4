package fzfyml

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Task struct {
	source          string
	sourceTransform string
	query           string
	variables       Variables
	binds           Binds
	preview         Preview
	options         Options
	postOperations  PostOperations
	switchExpects   []string
}

func (task *Task) init(baseTask map[string]interface{}, ymlPath string, switchExpects []string, args []string) {
	task.source = baseTask["source"].(string)
	if _, ok := baseTask["source_transform"]; ok {
		task.sourceTransform = baseTask["source_transform"].(string)
	}
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

func (task *Task) run(query interface{}) string {
	tmpTextName := ""
	tmpIndexName := ""
	defer os.Remove(tmpTextName)
	defer os.Remove(tmpIndexName)
	if task.sourceTransform != "" {
		tmpText, _ := ioutil.TempFile("", "fzfyml4-text-")
		tmpIndex, _ := ioutil.TempFile("", "fzfyml4-index-")
		tmpTextName = tmpText.Name()
		tmpIndexName = tmpIndex.Name()
	}
	command := task.getExecuteCommand("run", tmpTextName, tmpIndexName)
	if query != nil {
		command += " --query '" + query.(string) + "'"
	}
	resultText := task.execFzf(command)
	return resultText
}

func (task *Task) test(answer string) bool {
	tmpTextName := ""
	tmpIndexName := ""
	if task.sourceTransform != "" {
		tmpTextName = "./fzfyml4-text"
		tmpIndexName = "./fzfyml4-index"
	}
	response := task.getExecuteCommand("test", tmpTextName, tmpIndexName)
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

func (task *Task) getExecuteCommand(mode string, textFilePath string, indexFilePath string) string {
	source := task.getSourceText(textFilePath)
	bindList := task.binds.getBindList()
	preview := task.preview.getPreviewText(indexFilePath)
	optionList := task.options.getOptionList()
	expectList := task.getExpectList()
	mondatoryList := []string{"--print-query"}
	exe, _ := os.Executable()
	if mode == "test" {
		exe = "fzfyml4"
	}
	postCommand := task.getPostCommand(exe, textFilePath, indexFilePath)
	if mode == "test" {
		sort.Strings(bindList)
		sort.Strings(optionList)
		sort.Strings(expectList)
		sort.Strings(mondatoryList)
	}
	command := source + " | fzf " + strings.Join(bindList, " ") + " " + preview + " " + strings.Join(optionList, " ") + " --expect=" + strings.Join(expectList, ",") + " " + strings.Join(mondatoryList, " ")
	command = task.variables.expand(command)
	command += postCommand
	//fmt.Println(command + "\n")
	return command
}

func (task *Task) getPostCommand(exe string, textFilePath string, indexFilePath string) string {
	if task.sourceTransform == "" {
		return ""
	} else {
		return " | " + exe + " inner-untransformed-output " + textFilePath + " " + indexFilePath
	}
}

func (task *Task) getSourceText(textFilePath string) string {
	if task.sourceTransform == "" {
		return task.source
	} else {
		return task.source + " | tee " + textFilePath + " | " + task.sourceTransform
	}
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
