package main

import (
	"fmt"
	"isy-cli/internal/context"
	"isy-cli/internal/openai"
	"os"

	"github.com/spf13/cobra"
)

// ContextCommand definisce il comando CLI per generare il context
func ContextCommand() *cobra.Command {
	var verbose bool // Variabile per l'opzione verbose

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Genera un context unificato dai file specificati nel progetto",
		Run: func(cmd *cobra.Command, args []string) {
			outputPath := ".isy/last_context" // File dove salvare il contesto

			// Genera il contesto utilizzando BuildContext
			contextContent, err := context.BuildContext()
			if err != nil {
				fmt.Println("Errore durante la generazione del contesto:", err)
				return
			}

			// Se l'opzione verbose è attiva, stampa il contesto in console
			if verbose {
				fmt.Println("\n--- Context Content ---")
				fmt.Println(contextContent)
			}

			// Scrive il contesto su un file
			err = os.WriteFile(outputPath, []byte(contextContent), 0644)
			if err != nil {
				fmt.Println("Errore durante la scrittura del file del contesto:", err)
				return
			}

			fmt.Println("Contesto generato e salvato in:", outputPath)

			token, err := openai.TokenizerCtx(contextContent)
			if err != nil {
				fmt.Println("Errore durante la tokenizzazione del contesto:", err)
				return
			}

			fmt.Println("Numero di token della codebase totale:", len(token))

		},
	}

	// Aggiungi l'opzione verbose
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Stampa il contesto in console")

	return cmd
}
