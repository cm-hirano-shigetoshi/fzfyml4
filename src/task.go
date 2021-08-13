package fzfyml

import (
	"fmt"
	set "github.com/deckarep/golang-set"
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
	delimiter       interface{}
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
	if _, ok := baseTask["delimiter"]; ok {
		task.delimiter = baseTask["delimiter"].(string)
	}
	variables, _ := baseTask["variables"].(map[string]interface{})
	task.variables.init(ymlPath, args, variables)
	if _, ok := baseTask["binds"]; ok {
		task.binds.init(baseTask["binds"].(map[string]interface{}))
	} else {
		task.binds.init(nil)
	}
	if _, ok := baseTask["preview"]; ok {
		task.preview.init(baseTask["preview"].(map[string]interface{}))
	}
	if _, ok := baseTask["options"]; ok {
		task.options.init(baseTask["options"].([]interface{}))
	} else {
		task.options.init([]interface{}{})
	}
	task.options.setDelimiter(task.delimiter)
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
	command := task.getExecuteCommand("run", query, tmpTextName, tmpIndexName)
	resultText := task.execFzf(command)
	return resultText
}

func (task *Task) test(query interface{}, answer string) bool {
	tmpTextName := ""
	tmpIndexName := ""
	if task.sourceTransform != "" {
		tmpTextName = "./fzfyml4-text"
		tmpIndexName = "./fzfyml4-index"
	}
	command := task.getExecuteCommand("test", query, tmpTextName, tmpIndexName)
	if answer == command {
		return true
	} else {
		fmt.Println(answer)
		fmt.Println(command)
		return false
	}
}

func (task *Task) execFzf(command string) string {
	cmd_exec := exec.Command("bash", "-c", command)
	cmd_exec.Stderr = os.Stderr
	out, _ := cmd_exec.Output()
	if len(out) > 0 {
		out = out[:len(out)-1]
		return string(out)
	} else {
		return ""
	}
}

func (task *Task) getExecuteCommand(mode string, query interface{}, textFilePath string, indexFilePath string) string {
	exe, _ := os.Executable()
	if mode == "test" {
		exe = "fzfyml4"
	}
	source := task.getSourceText(textFilePath)
	bindList := task.binds.getBindList(exe, textFilePath, indexFilePath, task.delimiter)
	preview := task.preview.getPreviewText(exe, textFilePath, indexFilePath, task.delimiter)
	optionList := task.options.getOptionList()
	expectList := task.getExpectList()
	mondatoryList := []string{"--print-query"}
	queryCommand := ""
	if query != nil {
		queryCommand = "--query '" + query.(string) + "'"
	}
	postCommand := task.getPostCommand(exe, textFilePath, indexFilePath)
	if mode == "test" {
		sort.Strings(bindList)
		sort.Strings(optionList)
		sort.Strings(expectList)
		sort.Strings(mondatoryList)
	}
	command := concatStr(source, "|", "fzf", strings.Join(bindList, " "), preview, strings.Join(optionList, " "), "--expect="+strings.Join(expectList, ","), strings.Join(mondatoryList, " "), queryCommand)
	//fmt.Println(command)
	command = task.variables.expand(command)
	//fmt.Println(command)
	command = concatStr(command, postCommand)
	//fmt.Println(command + "\n")
	return command
}

func (task *Task) getPostCommand(exe string, textFilePath string, indexFilePath string) string {
	if task.sourceTransform == "" {
		return ""
	} else {
		return concatStr("|", exe, "inner-untransformed-output", textFilePath, indexFilePath)
	}
}

func (task *Task) getSourceText(textFilePath string) string {
	if task.sourceTransform == "" {
		return task.source
	} else {
		return concatStr(task.source, "|", "tee", textFilePath, "|", task.sourceTransform)
	}
}

func (task *Task) getExpectList() []string {
	expectSet := set.NewSet()
	for _, key := range task.postOperations.getExpects() {
		expectSet.Add(key)
	}
	for _, key := range task.switchExpects {
		expectSet.Add(key)
	}
	expectSet = expectSet.Union(set.NewSetFromSlice([]interface{}{"enter", "esc", "ctrl-c", "ctrl-d", "ctrl-g", "ctrl-q", "ctrl-z"}))
	for _, key := range task.binds.getBindKeys() {
		expectSet.Remove(key)
	}
	expects := []string{}
	for key := range expectSet.Iter() {
		expects = append(expects, key.(string))
	}
	return expects
}
