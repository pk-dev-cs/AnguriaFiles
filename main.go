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

	width  int
	height int
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

func (m *model) resizePicker() {
	if m.height <= 0 {
		return
	}

	reservedHeight :=
		lipgloss.Height(m.headerView()) +
			lipgloss.Height(m.footerView()) +
			2

	m.picker.SetHeight(max(m.height-reservedHeight, 1))
}

func (m model) headerView() string {
	width := m.width
	if width <= 0 {
		width = 80
	}

	return renderLogo(width) +
		"\n\n" +
		"  ────────────\n" +
		"  Katalog: " + m.picker.CurrentDirectory
}

func (m model) footerView() string {
	var content strings.Builder

	if m.selectedFile != "" {
		content.WriteString("  Wybrany plik: ")
		content.WriteString(m.selectedFile)
		content.WriteString("\n\n")
	}

	hiddenStatus := "wyłączone"
	if m.picker.ShowHidden {
		hiddenStatus = "włączone"
	}

	content.WriteString("  ↑/k, ↓/j     poruszanie\n")
	content.WriteString("  enter/l/→    otwórz katalog lub wybierz plik\n")
	content.WriteString("  h/←/esc      poprzedni katalog\n")
	content.WriteString("  .            pliki ukryte: ")
	content.WriteString(hiddenStatus)
	content.WriteString("\n")
	content.WriteString("  r            odśwież\n")
	content.WriteString("  q            wyjście")

	return content.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizePicker()

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
		m.resizePicker()
	}

	return m, cmd
}

func renderLogo(width int) string {
	const logo = `
⠀⠀⠀⠀⢀⣴⠂⣠⠖⠛⣿⣿⣿⡿⣿⣿⡟⢿⣿⡉⠻⣿⣿⠟⢁⡴⠃⡄⠀⠀
⠀⠀⠀⢠⡾⠁⣴⣷⣤⣾⣃⣾⡿⢁⣿⣿⣇⣀⣿⣿⡾⠟⢁⣴⠟⣡⣾⠁⠀⠀
⠀⠀⢠⡿⢁⣼⣿⣧⣴⣿⣿⣹⣿⣾⣿⣏⣻⡿⠟⢉⣠⡶⠟⣡⣾⠟⠁⣸⡇⠀
⠀⠀⣿⠇⣼⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠋⣁⣤⡶⠟⢉⣴⣾⠟⢁⠀⢰⣿⠀⠀
⠀⠸⣿⡄⢿⣿⣿⣿⡿⠿⠟⠋⣁⣤⡶⠿⠛⢁⣤⣾⡿⠋⢁⣴⠋⣠⠛⠋⠀⠀
⠀⠀⠻⢿⣦⣤⣤⣤⣤⠶⠞⠛⠉⡁⠄⠐⠚⠿⠛⢉⣠⡾⠟⢁⣴⠃⠈⠀⠀⠀
⠀⠀⠀⠠⠤⣬⣥⣤⣤⡴⠶⠟⢁⣤⣶⣷⠶⠒⠚⠋⣉⠤⠶⠿⠋⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠠⣤⣤⣤⣤⣴⣾⣿⣿⠟⢁⣴⣾⠟⣁⣤⣤⠖⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠉⠉⠛⠛⠛⢉⣠⣴⡿⠟⢁⡴⠟⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`

	// U+2800 wygląda jak spacja, ale nią nie jest.
	art := strings.ReplaceAll(logo, "\u2800", " ")

	// Usuwamy tylko początkowy i końcowy enter.
	art = strings.Trim(art, "\r\n")

	logoWidth := lipgloss.Width(art)

	title := lipgloss.NewStyle().
		Bold(true).
		Render("ANGURIA FILES")

	// Centrujemy tytuł względem logo.
	title = lipgloss.PlaceHorizontal(
		logoWidth,
		lipgloss.Center,
		title,
	)

	block := art + "\n\n" + title

	// Centrujemy cały blok, nie każdą linię osobno.
	return lipgloss.PlaceHorizontal(
		width,
		lipgloss.Left,
		block,
	)
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	content := strings.Join(
		[]string{
			m.headerView(),
			strings.TrimRight(m.picker.View(), "\r\n"),
			m.footerView(),
		}, "\n\n",
	)

	view := tea.NewView(content)
	view.BackgroundColor = lipgloss.Color("#0c1e3b")
	view.ForegroundColor = lipgloss.Color("#ffffff")
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
