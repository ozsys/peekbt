package main

import (
	"fmt"
	"os"

	"github.com/ozsys/peekbt/cmd/main/commands"
)

func hello() string {
	return "Welcome to peekbt!"
}

func goMain() int {
	if err := commands.Execute(); err != nil {
		fmt.Println(err) // Usageヘルプなどの出力
		return 1
	}
	return 0
}

func main() {
	status := goMain()
	os.Exit(status)
}
