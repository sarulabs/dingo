package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	mainFile := os.Getenv("GOPATH") + "/src/github.com/sarulabs/dingo/dingo/main.go"
	inputDir := os.Getenv("GOPATH") + "/src/github.com/sarulabs/dingo/dingo/tests/app/services"
	outputDir := os.Getenv("GOPATH") + "/src/github.com/sarulabs/dingo/dingo/tests/app/generated_services"

	out, err := exec.Command("go", "run", mainFile, "-src="+inputDir, "-dest="+outputDir).CombinedOutput()

	fmt.Println(string(out))

	if err != nil {
		os.Exit(1)
	}
}
