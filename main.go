package main

import (
	"dj/cmd"
	"fmt"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
