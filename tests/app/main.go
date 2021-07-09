package main

import (
	"fmt"
	"os"

	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/services/provider"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run main.go path/to/output/directory")
		os.Exit(1)
	}

	err := dingo.GenerateContainerWithCustomPkgName((*provider.Provider)(nil), os.Args[1], "dic")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
