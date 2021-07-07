package fzfyml

import ()

type Preview struct {
	command string
	window  string
}

func (preview *Preview) init(p map[string]interface{}) {
	preview.command = p["command"].(string)
	preview.window = p["window"].(string)
}

func (preview *Preview) getPreviewText() string {
	return "--preview '" + preview.command + "' --preview-window '" + preview.window + "'"
}
