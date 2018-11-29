package main

import (
	"fmt"
	"os"

	"github.com/izumin5210-sandbox/github-cli-sample/pkg/ghcp"
	"github.com/izumin5210-sandbox/github-cli-sample/pkg/ghcp/cmd"
)

func main() {
	var exitCode int

	if err := run(); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	os.Exit(exitCode)
}

func run() error {
	cmd := cmd.NewGhcpCommand(&ghcp.Ctx{
		IO: ghcp.StdIO(),
	})

	return cmd.Execute()
}
