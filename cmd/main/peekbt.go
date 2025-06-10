package main

import (
	"fmt"
	"os"
	"github.com/spf13/pflag"
)

func hello() string {
	return "Welcome to peekbt!"
}

func goMain(args []string) int {
	if err := commands.Execute(args); err != nil {
		fmt.Println(err) // Usageヘルプなどの出力
		return 1
	}
	fmt.Println(hello())
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}