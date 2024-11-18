package main

import (
	"isy-cli/internal/config"

	"github.com/spf13/cobra"
)

func InitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize the project",
		Run: func(cmd *cobra.Command, args []string) {
			config.InitProject()
		},
	}
}
