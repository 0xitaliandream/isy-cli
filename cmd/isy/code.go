package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	codeUtils "isy-cli/internal/code"
	"isy-cli/internal/context"
	localOpenAI "isy-cli/internal/openai"
	"isy-cli/internal/openai/schemas/code"
	"os"
	"path/filepath"
	"strconv"

	externalOpenAI "github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

func CodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "code",
		Short: "Interactively modify code with OpenAI assistance",
		Run: func(cmd *cobra.Command, args []string) {
			var selectedBranch string

			branchesDir := ".isy/branches"

			if len(args) > 0 {
				selectedBranch = args[0]
			} else {
				branches, err := codeUtils.ListBranchesSortedByDate(branchesDir)
				if err != nil {
					fmt.Println("Error retrieving branches:", err)
					return
				}

				if len(branches) == 0 {
				} else {
					fmt.Println("Please select a branch to connect:")
					for i, b := range branches[:min(3, len(branches))] {
						fmt.Printf("%d: %s\n", i+1, b)
					}
					fmt.Print("Enter the number of the branch to use or press Enter to create a new branch: ")

					var userChoice string
					fmt.Scanln(&userChoice)

					if choiceIndex, err := strconv.Atoi(userChoice); err == nil && choiceIndex >= 1 && choiceIndex <= min(3, len(branches)) {
						selectedBranch = branches[choiceIndex-1]
					}
				}
			}

			if selectedBranch == "" { // Nel caso non venga selezionato alcun branch esistente
				fmt.Println("No previous branch found, creating a new virtual branch.")
				hash, err := codeUtils.GenerateHash()
				if err != nil {
					fmt.Println("Error generating hash:", err)
					return
				}
				selectedBranch = hash
			}

			tempDir := filepath.Join(branchesDir, selectedBranch)

			if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
				currentHash, err := codeUtils.ComputeDirectoryHash(".")
				if err != nil {
					fmt.Println("Error computing current directory hash:", err)
					return
				}

				branchHash, err := codeUtils.ComputeDirectoryHash(tempDir)
				if err != nil {
					fmt.Println("Error computing branch directory hash:", err)
					return
				}

				if currentHash != branchHash {
					fmt.Println("The current codebase differs from the selected branch. Operation aborted.")
					return
				}
			} else {
				if err := codeUtils.CopyDir(".", tempDir); err != nil {
					fmt.Println("Error copying project:", err)
					return
				}
			}

			fmt.Println("Working on branch:", tempDir)

			// Inizializza il contesto
			contextContent, err := context.BuildContext()
			if err != nil {
				fmt.Println("Errore durante la generazione del contesto:", err)
				return
			}

			systemPrompt := code.SYSTEM_PROMPT

			// Initialize an empty chat session
			chat := []externalOpenAI.ChatCompletionMessageParamUnion{
				externalOpenAI.SystemMessage(systemPrompt),
				externalOpenAI.UserMessage(contextContent),
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Interactive code modification session started. Type your requests. Press Ctrl+D to exit.")

			for {
				fmt.Print("\nYou: ")
				userInput, err := reader.ReadString('\n')
				if err != nil {
					// EOF detected, exit the session
					fmt.Println("\nExiting chat session.")
					break
				}

				userInput = userInput[:len(userInput)-1] // Rimuovi il newline

				// Aggiungi il messaggio dell'utente alla chat
				chat = append(chat, externalOpenAI.UserMessage(userInput))

				// Prepare request parameters for OpenAI completion
				params := externalOpenAI.ChatCompletionNewParams{
					Model: externalOpenAI.F(externalOpenAI.ChatModelGPT4o),
					ResponseFormat: externalOpenAI.F[externalOpenAI.ChatCompletionNewParamsResponseFormatUnion](
						externalOpenAI.ResponseFormatJSONSchemaParam{
							Type: externalOpenAI.F(externalOpenAI.ResponseFormatJSONSchemaTypeJSONSchema),
							JSONSchema: externalOpenAI.F(externalOpenAI.ResponseFormatJSONSchemaJSONSchemaParam{
								Name:   externalOpenAI.F("code_modification"),
								Schema: externalOpenAI.F(code.CodeModificationResponseSchema),
								Strict: externalOpenAI.Bool(true),
							}),
						}),
					Messages: externalOpenAI.F(chat),
				}

				// Execute the completion request
				response, err := localOpenAI.RunCompletion(params)
				if err != nil {
					fmt.Println("Errore durante la richiesta a OpenAI:", err)
					continue
				}

				// Decode the response into the structured schema
				codeModificationResponse := code.CodeModificationResponse{}
				err = json.Unmarshal([]byte(response), &codeModificationResponse)
				if err != nil {
					fmt.Println("Errore nella decodifica della risposta:", err)
					continue
				}

				// Print the response
				fmt.Println("OpenAI:", codeModificationResponse)
			}
		},
	}
}
