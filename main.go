package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/filepicker"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	picker       filepicker.Model
	selectedFile string
	quitting     bool
}

func newModel() (model, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return model{}, fmt.Errorf(
			"nie udało się odczytać bieżącego katalogu: %w",
			err,
		)
	}

	picker := filepicker.New()

	picker.AutoHeight = false
	picker.CurrentDirectory = currentDirectory
	picker.ShowPermissions = true
	picker.ShowSize = true
	picker.ShowHidden = false

	// Katalogów nie wybieramy jako wynik.
	// Enter na katalogu otwiera katalog.
	picker.DirAllowed = false

	// Pliki mogą zostać wybrane.
	picker.FileAllowed = true

	return model{
		picker: picker,
	}, nil
}

func (m model) Init() tea.Cmd {
	return m.picker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		const reservedHeight = 16

		pickerHeight := max(msg.Height-reservedHeight, 1)

		m.picker.SetHeight(pickerHeight)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case ".":
			// Pokazywanie lub ukrywanie plików ukrytych.
			m.picker.ShowHidden = !m.picker.ShowHidden
			return m, m.picker.Init()

		case "r":
			// Ponowne odczytanie bieżącego katalogu.
			return m, m.picker.Init()
		}
	}

	var cmd tea.Cmd

	// Przekazujemy wiadomość do komponentu filepicker.
	m.picker, cmd = m.picker.Update(msg)

	// Sprawdzamy, czy użytkownik wybrał plik.
	if selected, path := m.picker.DidSelectFile(msg); selected {
		m.selectedFile = path
	}

	return m, cmd
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var content strings.Builder

	content.WriteString("\n")
	content.WriteString("  FILE MANAGER\n")
	content.WriteString("  ────────────\n")
	content.WriteString("  Katalog: ")
	content.WriteString(m.picker.CurrentDirectory)
	content.WriteString("\n\n")

	content.WriteString(m.picker.View())

	if m.selectedFile != "" {
		content.WriteString("\n")
		content.WriteString("  Wybrany plik: ")
		content.WriteString(m.selectedFile)
		content.WriteString("\n")
	}

	hiddenStatus := "wyłączone"
	if m.picker.ShowHidden {
		hiddenStatus = "włączone"
	}

	content.WriteString("\n")
	content.WriteString("  ↑/k, ↓/j     poruszanie\n")
	content.WriteString("  enter/l/→    otwórz katalog lub wybierz plik\n")
	content.WriteString("  h/←/esc      poprzedni katalog\n")
	content.WriteString("  .            pliki ukryte: ")
	content.WriteString(hiddenStatus)
	content.WriteString("\n")
	content.WriteString("  r            odśwież\n")
	content.WriteString("  q            wyjście\n")

	view := tea.NewView(content.String())
	view.AltScreen = true

	return view
}

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
