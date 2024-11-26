package context

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"isy-cli/internal/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/zabawaba99/go-gitignore"
)

func GenerateTree() (string, error) {
	baseDir := "."

	// Ottieni i file da includere dalla funzione GetFilesFromIsyContext
	filesFromContext, err := GetFilesFromIsyContext()
	if err != nil {
		return "", fmt.Errorf("errore durante il recupero dei file da .isycontext: %v", err)
	}

	// Usa una mappa per verifiche veloci
	fileSet := make(map[string]bool)
	for _, file := range filesFromContext {
		fileSet[file] = true
	}

	var treeBuilder strings.Builder

	var buildTree func(string, string) error
	buildTree = func(currentDir, prefix string) error {
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return fmt.Errorf("errore durante la lettura della directory %s: %v", currentDir, err)
		}

		var validEntries []os.DirEntry
		for _, entry := range entries {
			fullPath := filepath.Join(currentDir, entry.Name())
			// Includi solo i file e le directory che appartengono al contesto
			if entry.IsDir() {
				// Se una directory contiene file nel contesto, includila
				include := false
				for path := range fileSet {
					if strings.HasPrefix(path, fullPath) {
						include = true
						break
					}
				}
				if include {
					validEntries = append(validEntries, entry)
				}
			} else if fileSet[fullPath] {
				validEntries = append(validEntries, entry)
			}
		}

		for i, entry := range validEntries {
			isLast := i == len(validEntries)-1
			connector := "├── "
			if isLast {
				connector = "└── "
			}

			fullPath := filepath.Join(currentDir, entry.Name())
			treeBuilder.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, entry.Name()))

			if entry.IsDir() {
				newPrefix := prefix
				if isLast {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
				err := buildTree(fullPath, newPrefix)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Aggiungi la directory di base
	treeBuilder.WriteString(fmt.Sprintf("%s/\n", filepath.Base(baseDir)))
	err = buildTree(baseDir, "")
	if err != nil {
		return "", err
	}

	return treeBuilder.String(), nil
}

func MergeFiles() (string, error) {

	// Ottieni la lista dei file da GetFilesFromIsyContext
	filesFromContext, err := GetFilesFromIsyContext()
	if err != nil {
		return "", fmt.Errorf("errore durante il recupero dei file da .isycontext: %v", err)
	}

	var mergedContent strings.Builder

	for _, path := range filesFromContext {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue // Ignora file non validi o directory
		}

		// Leggi il file e aggiungilo al contenuto unificato
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("errore durante la lettura del file %s: %v", path, err)
		}

		lines := strings.Split(string(content), "\n")

		mergedContent.WriteString("----- START FILE -----\n")
		mergedContent.WriteString(fmt.Sprintf("FILE: %s\n", path))
		mergedContent.WriteString("----- CONTENT -----\n")

		for i, line := range lines {
			mergedContent.WriteString(fmt.Sprintf("%d: %s\n", i+1, line))
		}

		mergedContent.WriteString("----- END FILE -----\n\n")
	}

	return mergedContent.String(), nil
}

func readPatternsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("errore durante l'apertura del file: %v", err)
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Ignore comments and empty lines
		if line == "" || line[0] == '#' {
			continue
		}
		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("errore durante la lettura del file: %v", err)
	}

	return patterns, nil
}

func GetFilesFromIsyContext() ([]string, error) {
	baseDir := "."
	isyContextPath := filepath.Join(baseDir, ".isycontext")

	// Legge i pattern dal file
	filePatterns, err := readPatternsFromFile(isyContextPath)
	if err != nil {
		return nil, err
	}
	var matchedFiles []string

	// Scansiona i file nella directory
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Salta la directory ".isy"
		if info.IsDir() && info.Name() == ".isy" {
			return filepath.SkipDir
		}

		// Salta la directory ".isy"
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Salta directory
		if info.IsDir() {
			return nil
		}

		// Confronta i glob pattern
		for _, pattern := range filePatterns {
			// Usa il pattern per confrontare
			matched := gitignore.Match(pattern, path)
			if matched {
				matchedFiles = append(matchedFiles, path)
				break
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("errore durante il matching dei file: %v", err)
	}

	return matchedFiles, nil
}

func BuildContext() (string, error) {

	// Carica la configurazione del progetto
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("errore durante il caricamento della configurazione: %v", err)
	}

	// Genera l'albero dei file
	projectTree, err := GenerateTree()
	if err != nil {
		return "", fmt.Errorf("errore durante la generazione dell'albero del progetto: %v", err)
	}

	// Unisce i contenuti dei file specificati
	mergedContent, err := MergeFiles()
	if err != nil {
		return "", fmt.Errorf("errore durante la generazione del contesto dei file: %v", err)
	}

	// Costruisce il contesto come stringa
	var contextBuilder strings.Builder

	contextBuilder.WriteString("----- START CONTEXT -----\n\n")
	// Aggiungi informazioni sul progetto
	contextBuilder.WriteString("----- PROJECT INFO -----\n\n")
	contextBuilder.WriteString(fmt.Sprintf("Project Name: %s\n", cfg.ProjectName))
	contextBuilder.WriteString(fmt.Sprintf("Description: %s\n", cfg.Description))
	contextBuilder.WriteString("\n----- END PROJECT INFO -----\n\n")

	// Aggiungi l'albero dei file
	contextBuilder.WriteString("----- START PROJECT TREE -----\n\n")
	contextBuilder.WriteString(projectTree)
	contextBuilder.WriteString("\n----- END PROJECT TREE -----\n\n")

	// Aggiungi i contenuti unificati dei file
	contextBuilder.WriteString(mergedContent)
	contextBuilder.WriteString("----- END CONTEXT -----\n\n")

	return contextBuilder.String(), nil
}
