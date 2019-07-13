package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	mainFile := dir + "/../../main.go"
	inputDir := dir + "/services"
	outputDir := dir + "/generated_services"
	destPkg := "github.com/sarulabs/dingo/dingo/tests/app/generated_services"

	out, err := exec.Command("go", "run", mainFile, "-src="+inputDir, "-dest="+outputDir, "-destPkg="+destPkg).CombinedOutput()

	fmt.Println(string(out))

	if err != nil {
		os.Exit(1)
	}
}
