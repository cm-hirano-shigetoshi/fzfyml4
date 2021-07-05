package main

import (
	"flag"
	"log"
	"main/src"
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
	}
}
