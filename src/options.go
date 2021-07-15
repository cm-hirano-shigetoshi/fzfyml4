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

func (options *Options) setDelimiter(delimiter interface{}) {
	if delimiter != nil {
		options.list["delimiter"] = delimiter.(string)
	}
}

func (options *Options) getOptionList() []string {
	list := []string{}
	for key, val := range options.list {
		if val != "nil" {
			list = append(list, "--"+key+"='"+val+"'")
		} else {
			list = append(list, "--"+key)
		}
	}
	return list
}
