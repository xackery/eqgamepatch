package main

import (
	"fmt"
	"os"

	"github.com/xackery/eqgamepatch/ui"
)

func main() {

	ui, err := ui.New()
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
