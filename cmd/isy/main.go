package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "isy",
		Short: "isy is your AI-powered coding assistant",
		Long:  "isy is a CLI tool to help you manage and edit your codebase with the power of OpenAI.",
	}

	// Aggiungi i comandi disponibili
	rootCmd.AddCommand(InitCommand())
	rootCmd.AddCommand(AskCommand())
	rootCmd.AddCommand(CodeCommand())
	rootCmd.AddCommand(ContextCommand()) // Aggiunto il comando context

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
