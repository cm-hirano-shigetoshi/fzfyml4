package fzfyml

import (
	"strings"
	//"fmt"
)

type Options struct {
	list map[string]string
}

func (options *Options) InitFromYml(optionList []interface{}) {
	options.list = map[string]string{}
	options.setDefaultOptions()
	for _, opt := range optionList {
		if strings.Contains(opt.(string), "=") {
			sp := strings.Split(opt.(string), "=")
			options.list[sp[0]] = strings.Join(sp[1:], "=")
		} else {
			options.list[opt.(string)] = "nil"
		}
	}
}

func (options *Options) GetOptionText() string {
	opts := []string{}
	for key, val := range options.list {
		if val != "nil" {
			opts = append(opts, "--"+key+"="+val)
		} else {
			opts = append(opts, "--"+key)
		}
	}
	return strings.Join(opts, " ")
}

func (options *Options) setDefaultOptions() {
	options.list["print-query"] = "nil"
	options.list["expect"] = "enter"
}
