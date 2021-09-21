package fzfyml

import (
	"bufio"
	"fmt"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

func Run() int {
	delimiterPtr := flag.StringP("delimiter", "d", "__nil__", "delimiter")
	pathPtr := flag.StringP("path", "p", "auto", "auto|absolute|relative")
	updirPtr := flag.Int("updir_depth", 3, "auto時に絶対パスと相対パスを切り替える深さ")
	tildePtr := flag.Bool("tilde_home", false, "ホームディレクトリを~にする")
	slashPtr := flag.Bool("slash", false, "ディレクトリの場合末尾に/をつける")
	curdirPtr := flag.String("curdir", ".", "相対パスの起点パス")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("引数がありません")
		return 1
	}
	if flag.Args()[0] == "run" {
		execute(flag.Args()[1], flag.Args()[2:])
	} else if flag.Args()[0] == "test" {
		test(flag.Args()[1])
	} else if flag.Args()[0] == "inner-untransformed-output" {
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		fmt.Println(sc.Text())
		sc.Scan()
		fmt.Println(sc.Text())
		lines := selectLine(flag.Args()[1], flag.Args()[2])
		for _, line := range lines {
			fmt.Println(line)
		}
	} else if flag.Args()[0] == "nth-tool" {
		var delimiter interface{} = nil
		if *delimiterPtr != "__nil__" {
			delimiter = *delimiterPtr
		}
		fmt.Println(nth(flag.Args()[1], delimiter))
	} else if flag.Args()[0] == "path-tool" {
		fmt.Println(formatPath(*pathPtr, *updirPtr, *tildePtr, *slashPtr, *curdirPtr))
	}
	return 0
}
