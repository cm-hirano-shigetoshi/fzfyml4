package fzfyml

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func getIndexSlice(indexFilePath string) ([]int, int) {
	indexsString, err := ioutil.ReadFile(indexFilePath)
	if err != nil {
		log.Fatal(err)
	}
	indexes := []int{}
	max := 0
	for _, s := range strings.Split(strings.TrimRight(string(indexsString), "\n"), " ") {
		n, _ := strconv.Atoi(s)
		if n > max {
			max = n
		}
		indexes = append(indexes, n)
	}
	return indexes, max
}

func prepareLinesToPrint(textFilePath string, indexes []int, max int) map[int]string {
	indexSet := map[int]bool{}
	for _, i := range indexes {
		indexSet[i] = true
	}
	fh, _ := os.Open(textFilePath)
	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	lines := map[int]string{}
	i := 0
	for scanner.Scan() {
		if _, ok := indexSet[i]; ok {
			lines[i] = scanner.Text()
		}
		i++
		if i > max {
			break
		}
	}
	return lines
}

func printLines(textFilePath string, indexes []int, max int) {
	lines := prepareLinesToPrint(textFilePath, indexes, max)
	for _, n := range indexes {
		fmt.Println(lines[n])
	}
}

func SelectLine(textFilePath string, indexFilePath string) {
	indexes, max := getIndexSlice(indexFilePath)
	printLines(textFilePath, indexes, max)
}
