package main

import (
	"codetest/cmd"
	"fmt"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println("cmd 执行失败", err)
		return
	}
}
