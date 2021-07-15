package fzfyml

import (
	"strings"
)

type Binds struct {
	list map[string][]string
}

func (binds *Binds) init(list map[string]interface{}) {
	binds.list = map[string][]string{}
	binds.update(list)
}

func (binds *Binds) update(list map[string]interface{}) {
	for key, val := range list {
		operations := []string{}
		for _, ope := range val.([]interface{}) {
			operations = append(operations, ope.(string))
		}
		binds.list[key] = operations
	}
}

func (binds *Binds) getBindList(exe string, textFilePath string, indexFilePath string, delimiter interface{}) []string {
	list := []string{}
	for key, operations := range binds.list {
		opeList := []string{}
		for _, ope := range operations {
			if textFilePath == "" {
				opeList = append(opeList, ope)
			} else {
				opeList = append(opeList, expandFieldIndex(ope, exe, textFilePath, delimiter))
			}
		}
		list = append(list, "--bind '"+key+":"+strings.Join(opeList, "+")+"'")
	}
	return list
}
