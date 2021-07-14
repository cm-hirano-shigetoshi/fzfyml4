package fzfyml

import (
	//"fmt"
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func pos(index int, length int) int {
	if index >= 0 {
		return index
	} else {
		return length + index
	}
}

func idx(str [][]string, index int) []string {
	if index >= 0 {
		return str[index]
	} else {
		index = len(str) + index
		return str[index]
	}
}

func getRangeText(line string, start int, end int, d interface{}) string {
	pattern := regexp.MustCompile(`\S+\s*`)
	var islands []string
	if d == nil {
		islands = []string{}
		for _, m := range pattern.FindAllStringSubmatch(line, -1) {
			islands = append(islands, m[0])
		}
	} else {
		islands = strings.Split(line, d.(string))
	}
	if start > 0 {
		start--
	}
	if end == -1 {
		if d == nil {
			return strings.Join(islands[pos(start, len(islands)):], "")
		} else {
			return strings.Join(islands[pos(start, len(islands)):], d.(string))
		}
	} else {
		if end <= 0 {
			end++
		}
		if d == nil {
			return strings.Join(islands[pos(start, len(islands)):pos(end, len(islands))], "")
		} else {
			return strings.Join(islands[pos(start, len(islands)):pos(end, len(islands))], d.(string))
		}
	}
}

func oneRange(line string, i string, d interface{}) string {
	pattern := regexp.MustCompile(`\S+\s*`)
	if !strings.Contains(i, "..") {
		index, _ := strconv.Atoi(i)
		if index > 0 {
			index--
		}
		if d == nil {
			matches := pattern.FindAllStringSubmatch(line, -1)
			return idx(matches, index)[0]
		} else {
			return "aaa"
		}
	} else {
		if i == ".." {
			start := 1
			end := -1
			return getRangeText(line, start, end, nil)
		} else if strings.HasPrefix(i, "..") {
			start := 1
			end, _ := strconv.Atoi(i[2:])
			return getRangeText(line, start, end, nil)
		} else if strings.HasSuffix(i, "..") {
			start, _ := strconv.Atoi(i[:len(i)-2])
			end := -1
			return getRangeText(line, start, end, nil)
		} else {
			start, _ := strconv.Atoi(i[:strings.Index(i, ".")])
			end, _ := strconv.Atoi(i[strings.Index(i, ".")+2:])
			return getRangeText(line, start, end, nil)
		}
	}
}

func Nth(indexs string) string {
	sc := bufio.NewScanner(os.Stdin)
	text := []string{}
	for sc.Scan() {
		line := sc.Text()
		lineElements := []string{}
		if indexs == "" {
			lineElements = append(lineElements, line)
		} else {
			for _, index := range strings.Split(indexs, ",") {
				lineElements = append(lineElements, oneRange(line, index, nil))
			}
		}
		text = append(text, strings.Join(lineElements, " "))
	}
	return strings.Join(text, "\n")
}
