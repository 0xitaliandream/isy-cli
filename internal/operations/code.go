package operations

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ModifyFile(filePath string, startLine, endLine int, newCode string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	// Controllo che startLine ed endLine siano validi
	if startLine < 1 || endLine > len(lines) || startLine > endLine {
		fmt.Println("Invalid line range")
		return
	}

	// Sostituzione del codice
	modifiedLines := append(lines[:startLine-1], append([]string{newCode}, lines[endLine:]...)...)
	modifiedContent := strings.Join(modifiedLines, "\n")

	// Scrittura del nuovo contenuto nel file
	err = ioutil.WriteFile(filePath, []byte(modifiedContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File modified successfully!")
}
