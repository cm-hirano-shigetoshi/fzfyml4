package fzfyml

import (
	//"fmt"
	"strings"
)

type Preview struct {
	command string
	window  string
}

func (preview *Preview) init(p map[string]interface{}) {
	if command, ok := p["command"].(string); ok {
		preview.command = command
	}
	if window, ok := p["window"].(string); ok {
		preview.window = window
	}
}

func (preview *Preview) getPreviewText(exe string, textFilePath string, indexFilePath string, delimiter interface{}) string {
	options := []string{}
	if len(preview.command) > 0 {
		options = append(options, "--preview")
		command := ""
		if indexFilePath == "" {
			command = "'" + preview.command + "'"
		} else {
			expanded_command := expandFieldIndex(preview.command, exe, textFilePath, delimiter)
			command = "'" + "echo {+n} > " + indexFilePath + "; " + expanded_command + "'"
		}
		options = append(options, command)
	}
	if len(preview.window) > 0 {
		options = append(options, "--preview-window")
		options = append(options, preview.window)
	}
	return strings.Join(options, " ")
}
