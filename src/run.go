package fzfyml

import (
	"fmt"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	app = kingpin.New("fzfyml4", "top command")
)

func Run() int {
	cmdRun := app.Command("run", "run")
	cmdRunYml := cmdRun.Arg("yml", "yml").String()
	cmdRunArgs := cmdRun.Arg("args", "args").Strings()

	cmdTest := app.Command("test", "test")
	cmdTestYml := cmdTest.Arg("yml", "yml").String()

	cmdInnterUntransformedOutput := app.Command("inner-untransformed-output", "inner-untransformed-output")
	cmdInnterUntransformedOutputTextFilePath := cmdInnterUntransformedOutput.Arg("text-file-path", "text-file-path").String()
	cmdInnterUntransformedOutputIndexFilePath := cmdInnterUntransformedOutput.Arg("index-file-path", "index-file-path").String()

	cmdNth := app.Command("nth-tool", "nth-tool")
	cmdNthIndexes := cmdNth.Arg("indexes", "indexes").String()
	cmdNthDelimiter := cmdNth.Flag("delimiter", "delimiter").Default("__nil__").Short('d').String()

	cmdPath := app.Command("path-tool", "path-tool")
	cmdPathPath := cmdPath.Flag("path", "auto|absolute|relative").Default("auto").Short('p').String()
	cmdPathUpdir := cmdPath.Flag("updir_depth", "auto時に絶対パスと相対パスを切り替える深さ").Default("3").Int()
	cmdPathTilde := cmdPath.Flag("tilde_home", "ホームディレクトリを~にする").Default("false").Bool()
	cmdPathSlash := cmdPath.Flag("slash", "ディレクトリの場合末尾に/をつける").Default("false").Bool()
	cmdPathCurdir := cmdPath.Flag("curdir", "相対パスの起点パス").Default(".").String()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdRun.FullCommand():
		exit, err := execute(*cmdRunYml, *cmdRunArgs)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdTest.FullCommand():
		exit, err := test(*cmdTestYml)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdInnterUntransformedOutput.FullCommand():
		exit, err := innerUntransformedOutput(*cmdInnterUntransformedOutputTextFilePath, *cmdInnterUntransformedOutputIndexFilePath)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdNth.FullCommand():
		delimiter := interface{}(nil).(string)
		if *cmdNthDelimiter != "__nil__" {
			delimiter = *cmdNthDelimiter
		}
		exit, err := nth(*cmdNthIndexes, delimiter)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdPath.FullCommand():
		exit, err := path(*cmdPathPath, *cmdPathUpdir, *cmdPathTilde, *cmdPathSlash, *cmdPathCurdir)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	}
	return 1
}
