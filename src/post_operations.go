package fzfyml

import (
	"io"
	"log"
	"os/exec"
	"strings"
	//"fmt"
)

type PostOperations struct {
	list map[string]interface{}
}

func (operations *PostOperations) init(operationList map[string]interface{}) {
	operations.list = operationList
}

func (operations *PostOperations) apply(result Result) string {
	output := strings.Join(result.output, "\n")
	if ope, ok := operations.list[result.key].([]interface{}); ok {
		for _, operation := range ope {
			switch operation.(type) {
			case string:
				if operation.(string) == "join" {
					output = strings.Replace(output, "\n", " ", -1)
				}
			case map[string]interface{}:
				if pipe, ok := operation.(map[string]interface{})["pipe"].(string); ok {
					cmd_exec := exec.Command("sh", "-c", pipe)
					stdin, err := cmd_exec.StdinPipe()
					if err != nil {
						log.Fatal(err)
					}
					go func() {
						defer stdin.Close()
						io.WriteString(stdin, output)
					}()
					out, err := cmd_exec.CombinedOutput()
					if err != nil {
						log.Fatal(err)
					}
					output = string(out)
				}
			}
		}
	}
	return output
}
