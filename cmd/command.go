package cmd

import (
	"apiTools/utils"
	"fmt"
	"os"
)

var commandStr string

var commandList = []string{"all", "serve", "bankInfo", "proxyPool", "crontab", "token"}

var helpString = `apiTools command line tool

Usage:
  apiTools run [command]

Available Commands:
	all 		Start running all programs in one process.
	serve		Start running apiTools program.
	bankInfo 	When BankBinInfo data does not exist in the database, insert BankBinInfo data into the database.
	proxyPool	Start running proxy pool program.
	crontab 	Start running crontab program.
	token		Manage user access token.

Use "apiTools help" for more information about a command.
`

func commandHelp() {
	fmt.Println(helpString)
}

func parseCommand() bool {
	// 接受命令行参数
	commandSlice := os.Args
	commandLen := len(commandSlice)
	if commandLen <= 1 || commandLen > 3 {
		commandHelp()
		return false
	}
	if commandSlice[1] != "run" {
		commandHelp()
		return false
	}

	if !utils.IsInSelic(commandSlice[2], commandList) {
		commandHelp()
		return false
	}
	commandStr = commandSlice[2]

	return true
}

// 处理命令行输入参数
func command() (err error) {
	switch commandStr {
	case "all":
		err = runAllProgram()
	case "serve":
		err = runHttpServer()
	case "bankInfo":
		err = runBankInfoTask()
	case "proxyPool":
		err = runProxyPoolApp()
	case "crontab":
		err = runCrontabApp()
	case "token":
		err = runTokenApp()
	}
	return
}
