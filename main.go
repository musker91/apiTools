package main

import (
	"apiTools/cmd"
	"apiTools/modles"
)

func main() {
	err := cmd.InitApiTools()

	defer modles.CloseIO()

	if err != nil {
		panic(err)
	}
}
