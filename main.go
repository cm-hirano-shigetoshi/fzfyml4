package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"main/src"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("引数がありません")
	}
	if flag.Args()[0] == "run" {
		fzfyml.Run(flag.Args()[1], flag.Args()[2:])
	} else if flag.Args()[0] == "test" {
		fzfyml.Test(flag.Args()[1])
	} else if flag.Args()[0] == "inner-untransformed-output" {
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		fmt.Println(sc.Text())
		sc.Scan()
		fmt.Println(sc.Text())
		lines := fzfyml.SelectLine(flag.Args()[1], flag.Args()[2])
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}
