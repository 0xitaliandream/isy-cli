package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"isy-cli/internal/context"
	localOpenAI "isy-cli/internal/openai"
	"isy-cli/internal/openai/schemas/ask"
	"os"

	externalOpenAI "github.com/openai/openai-go"

	"github.com/spf13/cobra"
)

func AskCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ask",
		Short: "Start a chat session with OpenAI about your codebase",
		Long:  "Engage in an interactive chat session to ask OpenAI questions about your codebase.",
		Run: func(cmd *cobra.Command, args []string) {
			// Inizializza il contesto
			contextContent, err := context.BuildContext()
			if err != nil {
				fmt.Println("Errore durante la generazione del contesto:", err)
				return
			}

			// Carica l'uso dei token all'inizio della sessione
			initialUsage, err := localOpenAI.LoadTokenUsage()
			if err != nil {
				fmt.Println("Errore caricando uso token iniziale:", err)
				return
			}

			// Usa lo schema JSON e il prompt dal pacchetto schemas/ask
			systemPrompt := ask.SYSTEM_PROMPT

			// Inizializza la sessione di chat
			chat := []externalOpenAI.ChatCompletionMessageParamUnion{
				externalOpenAI.SystemMessage(systemPrompt),
				externalOpenAI.UserMessage(contextContent),
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Chat session started. Type your queries. Press Ctrl+D to exit.")

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

				// Prepara i parametri per la chiamata a OpenAI
				params := externalOpenAI.ChatCompletionNewParams{
					Model: externalOpenAI.F(externalOpenAI.ChatModelGPT4o),
					ResponseFormat: externalOpenAI.F[externalOpenAI.ChatCompletionNewParamsResponseFormatUnion](
						externalOpenAI.ResponseFormatJSONSchemaParam{
							Type: externalOpenAI.F(externalOpenAI.ResponseFormatJSONSchemaTypeJSONSchema),
							JSONSchema: externalOpenAI.F(externalOpenAI.ResponseFormatJSONSchemaJSONSchemaParam{
								Name:   externalOpenAI.F("ask_code_info"),
								Schema: externalOpenAI.F(ask.AskCodeInfoResponseSchema),
								Strict: externalOpenAI.Bool(true),
							}),
						},
					),
					Messages: externalOpenAI.F(chat),
				}

				// Esegui la richiesta di completamento
				response, err := localOpenAI.RunCompletion(params)
				if err != nil {
					fmt.Println("Errore durante la richiesta a OpenAI:", err)
					continue
				}

				// Decodifica la risposta
				askResponse := ask.AskCodeInfo{}
				_ = json.Unmarshal([]byte(response), &askResponse)

				// Aggiungi la risposta al contesto della chat
				chat = append(chat, externalOpenAI.AssistantMessage(askResponse.ContextualResponse))

				// Mostra la risposta
				fmt.Printf("\nAssistant: %s\n", askResponse.ContextualResponse)
			}

			// Carica l'uso dei token alla fine della sessione
			finalUsage, err := localOpenAI.LoadTokenUsage()
			if err != nil {
				fmt.Println("Errore caricando uso token finale:", err)
				return
			}

			// Calcola i token e i costi della sessione
			inputTokensUsed := finalUsage.TokenInput - initialUsage.TokenInput
			outputTokensUsed := finalUsage.TokenOutput - initialUsage.TokenOutput
			sessionCost := finalUsage.TotalCost - initialUsage.TotalCost

			// Mostra i dettagli della sessione
			fmt.Printf("\nToken di input usati nella sessione: %d\n", inputTokensUsed)
			fmt.Printf("Token di output usati nella sessione: %d\n", outputTokensUsed)
			fmt.Printf("Costo della sessione (USD): %.4f\n", sessionCost)

			// Mostra i dettagli totali
			fmt.Printf("\nTotale token di input: %d\n", finalUsage.TokenInput)
			fmt.Printf("Totale token di output: %d\n", finalUsage.TokenOutput)
			fmt.Printf("Costo totale (USD): %.4f\n", finalUsage.TotalCost)
		},
	}
}
