package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"isy-cli/internal/config"
	"os"
	"sync"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/pkoukk/tiktoken-go"
)

var mu sync.Mutex // Per gestire l'accesso concorrente ai token globali

type TokenUsage struct {
	TokenInput  int64   `json:"token_input"`
	TokenOutput int64   `json:"token_output"`
	TotalCost   float64 `json:"total_cost"`
}

// Salva l'utilizzo dei token su disco
func SaveTokenUsage(usage TokenUsage) error {
	filePath := ".isy/token_usage.json"
	data, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		return fmt.Errorf("errore durante la serializzazione dei dati token: %v", err)
	}
	return ioutil.WriteFile(filePath, data, 0644)
}

// Carica l'utilizzo dei token da disco
func LoadTokenUsage() (TokenUsage, error) {

	filePath := ".isy/token_usage.json"

	var usage TokenUsage
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Se il file non esiste, restituisce un utilizzo vuoto
		return usage, nil
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return usage, fmt.Errorf("errore durante la lettura del file token: %v", err)
	}
	err = json.Unmarshal(data, &usage)
	if err != nil {
		return usage, fmt.Errorf("errore durante il parsing del file token: %v", err)
	}
	return usage, nil
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func TokenizerCtx(ctx string) ([]int, error) {
	encoding := "gpt-4o"

	tkm, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		fmt.Println("Errore durante la creazione del tokenizzatore:", err)
		return nil, err
	}

	token := tkm.Encode(ctx, nil, nil)

	return token, nil
}

func RunCompletion(params openai.ChatCompletionNewParams) (string, error) {
	// Percorso al file di configurazione e al file di utilizzo dei token

	// Carica la configurazione
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("errore durante il caricamento della configurazione: %v", err)
	}

	// Carica l'utilizzo corrente dei token
	currentUsage, err := LoadTokenUsage()
	if err != nil {
		return "", fmt.Errorf("errore durante il caricamento dell'utilizzo dei token: %v", err)
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey), // Imposta la chiave API
	)

	ctx := context.Background()

	// Esegui la richiesta di completamento
	chat, err := client.Chat.Completions.New(ctx, params, option.WithMaxRetries(5))
	if err != nil {
		return "", fmt.Errorf("errore durante la richiesta di completamento: %v", err)
	}

	// Aggiorna i token globali e calcola i costi
	mu.Lock()
	currentUsage.TokenInput += chat.Usage.PromptTokens
	currentUsage.TokenOutput += chat.Usage.CompletionTokens

	inputCost := float64(chat.Usage.PromptTokens) / 1_000_000 * 2.50
	outputCost := float64(chat.Usage.CompletionTokens) / 1_000_000 * 10
	currentUsage.TotalCost += inputCost + outputCost

	err = SaveTokenUsage(currentUsage)
	mu.Unlock()
	if err != nil {
		return "", fmt.Errorf("errore durante il salvataggio dell'utilizzo dei token: %v", err)
	}

	// Ritorna direttamente la risposta grezza
	if len(chat.Choices) > 0 {
		return chat.Choices[0].Message.Content, nil
	}

	// Se non ci sono scelte nella risposta, ritorna stringa vuota e un errore
	return "", fmt.Errorf("nessuna risposta disponibile dalla completion")
}
