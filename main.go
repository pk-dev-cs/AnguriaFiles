package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	initialModel, err := newModel()
	if err != nil {
		fmt.Fprintln(os.Stderr, "błąd:", err)
		os.Exit(1)
	}

	program := tea.NewProgram(initialModel)

	finalModel, err := program.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "błąd aplikacji:", err)
		os.Exit(1)
	}

	result, ok := finalModel.(model)
	if !ok {
		return
	}

	if result.selectedFile != "" {
		fmt.Println("Wybrany plik:", result.selectedFile)
	}
}
