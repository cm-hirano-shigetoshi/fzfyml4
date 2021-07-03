package fzfyml

import (
	"strings"
)

type Result struct {
	query  string
	key    string
	output []string
}

func (result *Result) Init(text string) {
	sp := strings.Split(text, "\n")
	result.query = sp[0]
	result.key = sp[1]
	result.output = sp[2:]
}
