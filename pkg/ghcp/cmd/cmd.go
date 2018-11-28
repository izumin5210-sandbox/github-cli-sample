package cmd

import (
	"github.com/izumin5210/ghcp/pkg/ghcp"
	"github.com/spf13/cobra"
)

func NewGhcpCommand(ctx *ghcp.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ghcp",
		RunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
