package fzfyml

import (
	"strings"
)

type Binds struct {
	list map[string][]string
}

func (binds *Binds) init(list map[string]interface{}) {
	binds.list = map[string][]string{}
	for key, val := range list {
		for _, ope := range val.([]interface{}) {
			binds.list[key] = append(binds.list[key], ope.(string))
		}
	}
}

func (binds *Binds) getBindList() []string {
	list := []string{}
	for key, operations := range binds.list {
		opeList := []string{}
		for _, ope := range operations {
			opeList = append(opeList, ope)
		}
		list = append(list, "--bind '"+key+":"+strings.Join(opeList, "+")+"'")
	}
	return list
}
