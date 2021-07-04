package fzfyml

import (
	"strings"
	//"fmt"
)

type Options struct {
	list map[string]string
}

func (options *Options) init(optionList []interface{}) {
	options.list = map[string]string{}
	options.setDefaultOptions()
	options.initFromYml(optionList)
}

func (options *Options) initFromYml(optionList []interface{}) {
	for _, opt := range optionList {
		if strings.Contains(opt.(string), "=") {
			sp := strings.Split(opt.(string), "=")
			options.list[sp[0]] = strings.Join(sp[1:], "=")
		} else {
			options.list[opt.(string)] = "nil"
		}
	}
}

func (options *Options) getOptionText() string {
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
	options.list["expect"] = "enter,esc,ctrl-c,ctrl-d,ctrl-g,ctrl-q,ctrl-z"
}
