package fzfyml

import (
	//"fmt"
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type PathOption struct {
	path       string
	updirDepth int
	tildeHome  bool
	slash      bool
	curdir     string
}

func transform(line string, opt PathOption) string {
	isDir := false
	home, _ := os.LookupEnv("HOME")
	if line == "~" {
		line = home
		isDir = true
	}
	if strings.HasPrefix(line, "~/") {
		line = strings.Join([]string{home, line[2:]}, "/")
		fInfo, _ := os.Stat(line)
		isDir = fInfo.IsDir()
	}
	if opt.path == "absolute" {
		fInfo, _ := os.Stat(line)
		isDir = fInfo.IsDir()
		line, _ = filepath.Abs(line)
	} else if opt.path == "relative" {
		fInfo, _ := os.Stat(line)
		isDir = fInfo.IsDir()
		line = relativePath(opt.curdir, line)
	} else {
		fInfo, _ := os.Stat(line)
		isDir = fInfo.IsDir()
		line = relativePath(opt.curdir, line)
		up := ""
		for i := 0; i < opt.updirDepth; i++ {
			up += "../"
		}
		if strings.HasPrefix(line, up) {
			line, _ = filepath.Abs(line)
		}
	}
	if opt.slash && isDir {
		line += "/"
		line = strings.Replace(line, "//", "/", -1)
	}
	if opt.tildeHome && strings.HasPrefix(line, home) {
		line = "~" + line[len(home):]
	}
	return line
}

func Path(path string, updirDepth int, tildeHome bool, slash bool, curdir string) string {
	pathOption := PathOption{path, updirDepth, tildeHome, slash, curdir}
	sc := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for sc.Scan() {
		line := sc.Text()
		line = transform(line, pathOption)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
