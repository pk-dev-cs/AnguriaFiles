package main

import (
	"strings"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/lipgloss/v2"
)

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
