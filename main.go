package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"

	"palette/cmd"
)

var version = "0.3.0"

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(version),
	); err != nil {
		os.Exit(1)
	}
}
