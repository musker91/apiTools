package token

import (
	"apiTools/modles"
	"apiTools/utils"
	"fmt"
	"os"
)

var helpString = `user token manage tools

add 		Add user token, if exists update.
del		Delete user token.
flush		Flush user token to cache.
help 		Get user token manage tools help info.
exit 		Exit apiTools token app.

Use "help" for more information about a command.`

// 启动运行token管理
func RunTokenApp() {
	fmt.Println("use [help] cat help info.")
	for {
		var cmdStr string
		fmt.Print(">>> ")
		_, err := fmt.Scanf("%s", &cmdStr)
		if err != nil {
			fmt.Println("The input is incorrect, please re-enter")
			continue
		}
		switch cmdStr {
		case "add":
			addCmd()
		case "del":
			delCmd()
		case "flush":
			flushCmd()
		case "exit", "quit":
			os.Exit(0)
		default:
			helpCmd()
		}
	}
}

func helpCmd() {
	fmt.Println(helpString)
}

func addCmd() {
	var qq string
	fmt.Print("input user qq number: ")
	_, err := fmt.Scanf("%s", &qq)
	if err != nil {
		fmt.Println("user token create fail, please retry")
		return
	}
	token := utils.CreateToken()
	userToken := &modles.UserToken{
		QQ:    qq,
		Token: token,
	}
	err = modles.CreateToken(userToken)
	if err != nil {
		fmt.Println("user token create fail, please retry")
		return
	}
	flushCmd()
	fmt.Printf("token: %s\n", token)
	fmt.Println("create user token info success !!!")
	return
}

func delCmd() {
	var token string
	fmt.Print("input user token: ")
	_, err := fmt.Scanf("%s", &token)
	if err != nil {
		fmt.Println("user token delete fail, please retry")
		return
	}
	userToken := &modles.UserToken{
		Token: token,
	}
	err = modles.DeleteToken(userToken)
	if err != nil {
		fmt.Println("user token delete fail, please retry")
		return
	}
	flushCmd()
	fmt.Println("delete user token info success !!!")
	return
}

func flushCmd() {
	tokens, err := modles.GetTokensFromDB()
	if err != nil {
		fmt.Println("flush user token to db fail")
		return
	}
	var tokenSlice []string
	for _, t := range tokens {
		tokenSlice = append(tokenSlice, t.Token)
	}
	err = modles.WriteTokensToCache(tokenSlice)
	if err != nil {
		fmt.Println("flush user token to db fail")
		return
	}
	fmt.Println("flush user tokens success !!!")
	return
}
