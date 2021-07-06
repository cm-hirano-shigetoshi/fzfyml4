package fzfyml

import (
	"io"
	"log"
	"os/exec"
)

func execCommandWithInput(command string, input string) string {
	cmd_exec := exec.Command("sh", "-c", command)
	stdin, err := cmd_exec.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()
	out, err := cmd_exec.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	if len(out) > 0 {
		return string(out[:len(out)-1])
	} else {
		return ""
	}
}
