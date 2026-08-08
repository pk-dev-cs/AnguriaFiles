package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	overlay "github.com/madicen/bubble-overlay"
)

func (m model) headerView() string {
	width := m.width
	if width <= 0 {
		width = 80
	}

	return renderLogo(width)
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

func popup(message string, width, height int) tea.View {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 3).
		Render(message)

	return tea.NewView(
		lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			box,
		),
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

	if m.popupMessage != "" {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 3).
			Render(m.popupMessage)

		content = overlay.OverlayViewInCenter(
			content,
			box,
			m.width,
			m.height,
		)
	}

	view := tea.NewView(content)
	view.BackgroundColor = lipgloss.Color("#0c1e3b")
	view.ForegroundColor = lipgloss.Color("#ffffff")
	view.AltScreen = true

	return view
}
