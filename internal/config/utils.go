package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config rappresenta la struttura del file di configurazione
type Config struct {
	ProjectName          string   `json:"project_name"`
	Author               string   `json:"author"`
	LanguageAndFramework string   `json:"language_and_framework"`
	Description          string   `json:"description"`
	Files                []string `json:"files"`
	APIKey               string   `json:"api_key"` // Nuovo campo
}

// LoadConfig legge il file di configurazione e restituisce un oggetto Config
func LoadConfig() (*Config, error) {
	configPath := ".isy/config.json"
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("errore durante l'apertura del file di configurazione: %v", err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("errore durante la decodifica del file di configurazione: %v", err)
	}
	return &config, nil
}

// SaveConfig salva un oggetto Config nel file specificato
func SaveConfig(config *Config) error {
	configPath := ".isy/config.json"
	configFile, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del file di configurazione: %v", err)
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ") // Formatta il JSON con indentazione
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("errore durante la scrittura del file di configurazione: %v", err)
	}
	return nil
}

// UpdateConfig aggiorna uno o più campi della configurazione e salva i cambiamenti
func UpdateConfig(updates map[string]interface{}) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Applica gli aggiornamenti
	for key, value := range updates {
		switch key {
		case "project_name":
			if v, ok := value.(string); ok {
				config.ProjectName = v
			}
		case "author":
			if v, ok := value.(string); ok {
				config.Author = v
			}
		case "language_and_framework":
			if v, ok := value.(string); ok {
				config.LanguageAndFramework = v
			}
		case "description":
			if v, ok := value.(string); ok {
				config.Description = v
			}
		case "files":
			if v, ok := value.([]string); ok {
				config.Files = v
			}
		}
	}

	// Salva le modifiche
	return SaveConfig(config)
}
