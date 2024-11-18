package context

import (
	"fmt"
	"io/ioutil"
	"isy-cli/internal/config"
	"os"
	"path/filepath"
	"strings"
)

func GenerateTree(baseDir string) (string, error) {
	var treeBuilder strings.Builder

	var buildTree func(string, string) error
	buildTree = func(currentDir, prefix string) error {
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return fmt.Errorf("errore durante la lettura della directory %s: %v", currentDir, err)
		}

		for i, entry := range entries {
			isLast := i == len(entries)-1
			connector := "├── "
			if isLast {
				connector = "└── "
			}

			treeBuilder.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, entry.Name()))

			if entry.IsDir() {
				newPrefix := prefix
				if isLast {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
				err := buildTree(filepath.Join(currentDir, entry.Name()), newPrefix)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Aggiungi la directory di base
	treeBuilder.WriteString(fmt.Sprintf("%s/\n", filepath.Base(baseDir)))
	err := buildTree(baseDir, "")
	if err != nil {
		return "", err
	}

	return treeBuilder.String(), nil
}

// MergeFiles legge e concatena i file specificati da un array di estensioni
func MergeFiles(baseDir string, extensions []string) (string, error) {
	var mergedContent strings.Builder

	// Itera attraverso i file nella directory di base
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Controlla che il file abbia una delle estensioni specificate
		if !info.IsDir() {
			for _, ext := range extensions {
				if strings.HasSuffix(info.Name(), ext) {
					// Legge il contenuto del file
					content, err := ioutil.ReadFile(path)
					if err != nil {
						return fmt.Errorf("errore durante la lettura del file %s: %v", path, err)
					}

					// Aggiunge il contenuto al risultato con i delimitatori
					mergedContent.WriteString("----- START FILE -----\n")
					mergedContent.WriteString(fmt.Sprintf("FILE: %s\n", path))
					mergedContent.WriteString("----- CONTENT -----\n")
					mergedContent.WriteString(string(content))
					mergedContent.WriteString("\n----- END FILE -----\n\n")
				}
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("errore durante l'iterazione dei file: %v", err)
	}

	return mergedContent.String(), nil
}

func BuildContext() (string, error) {

	// Configurazione di base
	baseDir := "." // Directory di base

	// Carica la configurazione del progetto
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("errore durante il caricamento della configurazione: %v", err)
	}

	// Genera l'albero dei file
	projectTree, err := GenerateTree(baseDir)
	if err != nil {
		return "", fmt.Errorf("errore durante la generazione dell'albero del progetto: %v", err)
	}

	// Unisce i contenuti dei file specificati
	mergedContent, err := MergeFiles(baseDir, cfg.Files)
	if err != nil {
		return "", fmt.Errorf("errore durante la generazione del contesto dei file: %v", err)
	}

	// Costruisce il contesto come stringa
	var contextBuilder strings.Builder

	// Aggiungi informazioni sul progetto
	contextBuilder.WriteString("----- PROJECT INFO -----\n")
	contextBuilder.WriteString(fmt.Sprintf("Project Name: %s\n", cfg.ProjectName))
	contextBuilder.WriteString(fmt.Sprintf("Description: %s\n", cfg.Description))
	contextBuilder.WriteString("----- END PROJECT INFO -----\n\n")

	// Aggiungi l'albero dei file
	contextBuilder.WriteString("----- START PROJECT TREE -----\n")
	contextBuilder.WriteString(projectTree)
	contextBuilder.WriteString("----- END PROJECT TREE -----\n\n")

	// Aggiungi i contenuti unificati dei file
	contextBuilder.WriteString("----- START CONTEXT -----\n")
	contextBuilder.WriteString(mergedContent)
	contextBuilder.WriteString("----- END CONTEXT -----\n\n")

	return contextBuilder.String(), nil
}
