package fzfyml

import (
	//"fmt"
)

type Preview struct {
	command string
	window  string
}

func (preview *Preview) init(p map[string]interface{}) {
	preview.command = p["command"].(string)
	preview.window = p["window"].(string)
}

func (preview *Preview) getPreviewText(exe string, textFilePath string, indexFilePath string) string {
	if indexFilePath == "" {
		return "--preview '" + preview.command + "' --preview-window '" + preview.window + "'"
	} else {
		expanded_command := expandFieldIndex(preview.command, exe, textFilePath)
		return "--preview '" + "echo {+n} > " + indexFilePath + "; " + expanded_command + "' --preview-window '" + preview.window + "'"
	}
}

func expandFieldIndex(command string, exe string, textFilePath string) string {
	targets := getReplaceTargets(command)
	for i := len(targets) - 1; i >= 0; i-- {
		start, end := targets[i][0], targets[i][1]
		index := command[start+1 : end-1]
		command = command[:start] + "$(sed -n $(({n}+1))p " + textFilePath + " | " + exe + " inner-nth \"" + index + "\")" + command[end:]
	}
	return command
}
