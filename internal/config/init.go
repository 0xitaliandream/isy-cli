package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InitProject() {

	// Verifica se la directory .isy esiste
	if _, err := os.Stat(".isy"); err == nil {
		fmt.Println("Il progetto è già stato inizializzato.")
		fmt.Print("Vuoi cancellare tutto e rifare l'init? (s/n): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "s" {
			fmt.Println("Operazione annullata.")
			return
		}

		// Cancella la directory
		err := os.RemoveAll(".isy")
		if err != nil {
			fmt.Println("Errore durante la cancellazione della directory:", err)
			return
		}
	}

	// Crea la directory .isy
	err := os.Mkdir(".isy", 0755)
	if err != nil {
		fmt.Println("Errore durante la creazione della directory:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Raccoglie i dati di configurazione dall'utente
	fmt.Println("Inserisci il nome del progetto:")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	fmt.Println("Inserisci il nome dell'autore o dell'azienda:")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSpace(author)

	fmt.Println("Inserisci il linguaggio e/o framework da utilizzare:")
	languageAndFramework, _ := reader.ReadString('\n')
	languageAndFramework = strings.TrimSpace(languageAndFramework)

	fmt.Println("Descrivi il progetto (obiettivi, funzionalità, ecc.):")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Println("Inserisci le estensioni dei file da includere (separate da spazi):")
	extensionsInput, _ := reader.ReadString('\n')
	extensions := strings.Fields(strings.TrimSpace(extensionsInput))

	fmt.Println("Inserisci la tua OpenAI API Key:")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	// Aggiungi l'API key alla configurazione
	config := &Config{
		ProjectName:          projectName,
		Author:               author,
		LanguageAndFramework: languageAndFramework,
		Description:          description,
		Files:                extensions,
		APIKey:               apiKey, // Salva l'API key
	}

	// Salva la configurazione utilizzando SaveConfig
	err = SaveConfig(config)
	if err != nil {
		fmt.Println("Errore durante il salvataggio della configurazione:", err)
		return
	}

	fmt.Println("Progetto inizializzato con successo!")
	fmt.Printf("Nome progetto: %s\n", projectName)
	fmt.Printf("Autore: %s\n", author)
	fmt.Printf("Linguaggio/Framework: %s\n", languageAndFramework)
	fmt.Printf("Descrizione: %s\n", description)
	fmt.Println("Estensioni configurate:", strings.Join(extensions, ", "))
}
