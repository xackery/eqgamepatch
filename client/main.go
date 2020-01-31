package main

import (
	"fmt"
	"os"

	"github.com/xackery/eqgamepatch/client/ui"
)

const (
	// Title to display on patcher
	Title = ""
)

func main() {
	ui, err := ui.New(Title)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ui.Start()
	if err != nil {
		fmt.Println("start", err)
		os.Exit(1)
	}
}
