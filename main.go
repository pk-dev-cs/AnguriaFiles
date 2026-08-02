package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/filepicker"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type pane = int

const (
	leftPane = iota
	rightPane
)

type model struct {
	leftPicker  filepicker.Model
	rightPicker filepicker.Model

	activePane   pane
	selectedFile string
	quitting     bool

	width  int
	height int
}

func newPicker(directory string) filepicker.Model {
	picker := filepicker.New()

	picker.AutoHeight = false
	picker.CurrentDirectory = directory
	picker.ShowPermissions = true
	picker.ShowSize = true
	picker.ShowHidden = false
	picker.DirAllowed = false
	picker.FileAllowed = true

	return picker
}

func newModel() (model, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return model{}, fmt.Errorf(
			"nie udało się odczytać bieżącego katalogu: %w",
			err,
		)
	}

	return model{
		leftPicker:  newPicker(currentDirectory),
		rightPicker: newPicker(currentDirectory),
		activePane:  leftPane,
	}, nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.leftPicker.Init(),
		m.rightPicker.Init(),
	)
}

func (m *model) activePicker() *filepicker.Model {
	if m.activePane == leftPane {
		return &m.leftPicker
	}

	return &m.rightPicker
}

func (m *model) resizePickers() {
	if m.height <= 0 {
		return
	}

	const (
		sectionGaps       = 2
		panelChromeHeight = 5
	)

	availableHeight :=
		m.height -
			lipgloss.Height(m.headerView()) -
			lipgloss.Height(m.footerView()) -
			sectionGaps -
			panelChromeHeight

	pickerHeight := max(availableHeight, 1)

	m.leftPicker.SetHeight(pickerHeight)
	m.rightPicker.SetHeight(pickerHeight)
}

func (m model) headerView() string {
	width := m.width
	if width <= 0 {
		width = 80
	}

	return renderLogo(width)
}

func renderPanel(
	title string,
	picker filepicker.Model,
	active bool,
	outerWidth int,
) string {
	borderColor := lipgloss.Color("#53657d")
	titleColor := lipgloss.Color("#9aa8bd")

	if active {
		borderColor = lipgloss.Color("#f6c453")
		titleColor = lipgloss.Color("#f6c453")
	}

	titleView := lipgloss.NewStyle().
		Bold(true).
		Foreground(titleColor).
		Render(title)

	content := strings.Join(
		[]string{
			titleView,
			picker.CurrentDirectory,
			"",
			picker.View(),
		},
		"\n",
	)

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	contentWidth := max(
		outerWidth-panelStyle.GetHorizontalFrameSize(),
		1,
	)

	return panelStyle.
		Width(contentWidth).
		MaxWidth(contentWidth).
		Render(content)
}

func (m model) splitView() string {
	const gap = 1

	width := m.width
	if width <= 0 {
		width = 80
	}

	availableWidth := max(width-gap, 2)

	leftWidth := availableWidth / 2
	rightWidth := availableWidth - leftWidth

	left := renderPanel(
		"LEWY PANEL",
		m.leftPicker,
		m.activePane == leftPane,
		leftWidth,
	)

	right := renderPanel(
		"PRAWY PANEL",
		m.rightPicker,
		m.activePane == rightPane,
		rightWidth,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		strings.Repeat(" ", gap),
		right,
	)
}

func (m model) footerView() string {
	var content strings.Builder

	activePicker := m.leftPicker
	activeName := "lewy"

	if m.activePane == rightPane {
		activePicker = m.rightPicker
		activeName = "prawy"
	}

	if m.selectedFile != "" {
		content.WriteString("  Wybrany plik: ")
		content.WriteString(m.selectedFile)
		content.WriteString("\n\n")
	}

	hiddenStatus := "wyłączone"
	if activePicker.ShowHidden {
		hiddenStatus = "włączone"
	}

	content.WriteString("  Aktywny panel: ")
	content.WriteString(activeName)
	content.WriteString("\n")
	content.WriteString("  tab          zmień panel\n")
	content.WriteString("  ↑/k, ↓/j     poruszanie\n")
	content.WriteString("  enter/l/→    otwórz katalog lub wybierz plik\n")
	content.WriteString("  h/←/esc      poprzedni katalog\n")
	content.WriteString("  .            pliki ukryte: ")
	content.WriteString(hiddenStatus)
	content.WriteString("\n")
	content.WriteString("  r            odśwież aktywny panel\n")
	content.WriteString("  q            wyjście")

	return content.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typedMsg.Width
		m.height = typedMsg.Height
		m.resizePickers()

	case tea.KeyPressMsg:
		switch typedMsg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			if m.activePane == leftPane {
				m.activePane = rightPane
			} else {
				m.activePane = leftPane
			}

			return m, nil

		case ".":
			picker := m.activePicker()
			picker.ShowHidden = !picker.ShowHidden

			return m, picker.Init()

		case "r":
			return m, m.activePicker().Init()
		}

		picker := m.activePicker()

		updatedPicker, cmd := picker.Update(msg)
		*picker = updatedPicker

		if selected, path := picker.DidSelectFile(msg); selected {
			m.selectedFile = path
			m.resizePickers()
		}

		return m, cmd
	}

	var leftCmd tea.Cmd
	var rightCmd tea.Cmd

	m.leftPicker, leftCmd = m.leftPicker.Update(msg)
	m.rightPicker, rightCmd = m.rightPicker.Update(msg)

	return m, tea.Batch(leftCmd, rightCmd)
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

	art := strings.ReplaceAll(logo, "\u2800", " ")

	art = strings.Trim(art, "\r\n")

	logoWidth := lipgloss.Width(art)

	title := lipgloss.NewStyle().
		Bold(true).
		Render("ANGURIA FILES")

	title = lipgloss.PlaceHorizontal(
		logoWidth,
		lipgloss.Center,
		title,
	)

	block := art + "\n\n" + title

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
			m.splitView(),
			m.footerView(),
		},
		"\n\n",
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
