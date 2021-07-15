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

func (preview *Preview) getPreviewText(exe string, textFilePath string, indexFilePath string, delimiter interface{}) string {
	if indexFilePath == "" {
		return "--preview '" + preview.command + "' --preview-window '" + preview.window + "'"
	} else {
		expanded_command := expandFieldIndex(preview.command, exe, textFilePath, delimiter)
		return "--preview '" + "echo {+n} > " + indexFilePath + "; " + expanded_command + "' --preview-window '" + preview.window + "'"
	}
}
