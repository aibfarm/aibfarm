package main

import (
	"fmt"
	"log"
	"os"

	ishell "github.com/abiosoft/ishell"
	web "github.com/aibfarm/aibfarm/web"

	// web "github.com/aibfarm/aibfarm/web"
	config "github.com/aibfarm/aibfarm/config"
)

var GitCommit string

func main() {
	// var str string
	// var err error

	// conf.SaveLogToDisk()
	log.SetFlags(log.LstdFlags | log.Llongfile)

	if len(GitCommit) > 8 {
		log.Printf("GIT Version:%s", GitCommit[0:8])
		// conf.GitCommit = GitCommit
	}

	ip := "local"

	welcome := fmt.Sprintf("AIB FARM [%s] ", ip)
	log.Printf("\x1b[31;1m  START: %s  \x1b[0m \n", welcome)

	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {

		case "test":
			web.Test()
		case "sample":
			config.SampleConfig()

		default:

		}

	} else {
		log.Printf("no param: aibfarm")
	}

	//!-- 初始测试设置
	// config.SampleConfig()

	iShell()
}

func iShell() {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()
	shell.SetHomeHistoryPath(".ishell_history")
	// display welcome info.
	shell.Println("AIB Farm Interactive Shell")
	shell.Run()
}
