package main

import (
	"fmt"
	"log"
	"os"

	ishell "github.com/abiosoft/ishell"
	web "github.com/aibfarm/web"

	// web "github.com/aibfarm/aibfarm/web"
	config "github.com/aibfarm/config"
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

	// config.AIBFarmConfig.API_ENDPOINT = "ws://192.168.2.138:21779" //TODO
	config.AIBFarmConfig.API_ENDPOINT = "wss://farm-api.8tc.ca" //TODO

	welcome := fmt.Sprintf("AIB FARM [%s] ", ip)
	log.Printf("\x1b[31;1m  START: %s  \x1b[0m \n", welcome)

	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {

		case "test":
			web.Start_WebServer()

		case "web":
			web.Start_WebServer()

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
