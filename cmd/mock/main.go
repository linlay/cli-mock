package main

import (
	"os"

	"github.com/linlay/cli-mock/internal/app"
)

func main() {
	os.Exit(app.Execute(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
