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

	fmt.Println("Inserisci la tua OpenAI API Key:")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	fmt.Println("Inserisci la lingua con cui isy ti risponderà (it, en, ecc):")
	language, _ := reader.ReadString('\n')
	language = strings.TrimSpace(language)

	// Aggiungi l'API key alla configurazione
	config := &Config{
		ProjectName:             projectName,
		Author:                  author,
		LanguageAndFramework:    languageAndFramework,
		Description:             description,
		APIKey:                  apiKey, // Salva l'API key
		IaModelResponseLanguage: language,
	}

	// Salva la configurazione utilizzando SaveConfig
	err = SaveConfig(config)
	if err != nil {
		fmt.Println("Errore durante il salvataggio della configurazione:", err)
		return
	}

	fmt.Print("Creazione del file .isycontext per specificare i file da includere nel contesto")
	err = os.WriteFile(".isycontext", []byte("# Add your file patterns here\n"), 0644)
	if err != nil {
		fmt.Println("Errore durante la creazione del file .isycontext:", err)
		return
	}
	fmt.Println("File .isycontext creato con successo. Puoi modificare il file per aggiungere le regex dei file da includere.")

	fmt.Println("Progetto inizializzato con successo!")
	fmt.Printf("Nome progetto: %s\n", projectName)
	fmt.Printf("Autore: %s\n", author)
	fmt.Printf("Linguaggio/Framework: %s\n", languageAndFramework)
	fmt.Printf("Descrizione: %s\n", description)
}
