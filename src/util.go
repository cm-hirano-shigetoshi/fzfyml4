package fzfyml

import (
	"io"
	"log"
	"os/exec"
	"regexp"
)

func execCommandWithInput(command string, input string) string {
	cmd_exec := exec.Command("bash", "-c", command)
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

func uniqueStringSlice(slice []string) []string {
	m := map[string]string{}
	for _, s := range slice {
		m[s] = ""
	}
	ret := []string{}
	for key, _ := range m {
		ret = append(ret, key)
	}
	return ret
}

func getReplaceTargets(s string) [][]int {
	pattern := regexp.MustCompile(`{(?:-?\d*|-?\d*\.\.\-?\d*)}`)
	matches := pattern.FindAllStringSubmatchIndex(s, -1)
	return matches
}
