package main

import (
	"fmt"
	"os"

	"github.com/xackery/eqgamepatch/s3d"
)

func main() {
	a, err := s3d.New("s3d/abysmal_obj.s3d")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if err := a.ExtractAll("s3d/out"); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

}
