package main

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func colorMarkedLine(
	line string,
	startMarker string,
	endMarker string,
	defaultStyle lipgloss.Style,
	markedStyle lipgloss.Style,
) string {
	var result strings.Builder

	for {
		start := strings.Index(line, startMarker)

		if start == -1 {
			if line != "" {
				result.WriteString(defaultStyle.Render(line))
			}
			break
		}

		if start > 0 {
			result.WriteString(defaultStyle.Render(line[:start]))
		}

		line = line[start+len(startMarker):]

		end := strings.Index(line, endMarker)
		if end == -1 {
			result.WriteString(defaultStyle.Render(line))
			break
		}

		if end > 0 {
			result.WriteString(markedStyle.Render(line[:end]))
		}

		line = line[end+len(endMarker):]
	}

	return result.String()
}

func colorMarkedText(
	text string,
	startMarker string,
	endMarker string,
	defaultStyle lipgloss.Style,
	markedStyle lipgloss.Style,
) string {
	lines := strings.Split(text, "\n")

	for index, line := range lines {
		lines[index] = colorMarkedLine(
			line,
			startMarker,
			endMarker,
			defaultStyle,
			markedStyle,
		)
	}

	return strings.Join(lines, "\n")
}

func renderLogo(width int) string {
	const logo = `
⠀⠀⠀⠀⢀⣴⠂[[⣠⠖⠛⣿⣿⣿⡿⣿⣿⡟⢿⣿⡉⠻⣿⣿⠟]]⢁⡴⠃⡄⠀⠀
⠀⠀⠀⢠⡾⠁[[⣴⣷⣤⣾⣃⣾⡿⢁⣿⣿⣇⣀⣿⣿⡾⠟]]⢁⣴⠟⣡⣾⠁⠀⠀
⠀⠀⢠⡿⠁[[⣼⣿⣧⣴⣿⣿⣹⣿⣾⣿⣏⣻⡿⠟]]⢉⣠⡶⠟⣡⣾⠟⠁⣸⡇⠀
⠀⠀⣿⠇[[⣼⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠋]]⣁⣤⡶⠟⢉⣴⣾⠟⢁⠀⢰⣿⠀⠀
⠀⠸⣿⡄[[⢿⣿⣿⣿⡿⠿⠟⠋]]⣁⣤⡶⠿⠛⢁⣤⣾⡿⠋⢁⣴⠋⣠⠛⠋⠀⠀
⠀⠀⠻⢿⣦⣤⣤⣤⣤⠶⠞⠛⠉⡁⠄⠐⠚⠿⠛⢉⣠⡾⠟⢁⣴⠃⠈⠀⠀⠀
⠀⠀⠀⠠⠤⣬⣥⣤⣤⡴⠶⠟⢁⣤⣶⣷⠶⠒⠚⠋⣉⠤⠶⠿⠋⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠠⣤⣤⣤⣤⣴⣾⣿⣿⠟⢁⣴⣾⠟⣁⣤⣤⠖⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠉⠉⠛⠛⠛⢉⣠⣴⡿⠟⢁⡴⠟⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`

	art := strings.ReplaceAll(logo, "\u2800", " ")
	art = strings.Trim(art, "\r\n")

	rindStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#55b86a")).
		Bold(true)

	fleshStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff4f64")).
		Bold(true)

	art = colorMarkedText(
		art,
		"[[",
		"]]",
		rindStyle,
		fleshStyle,
	)

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
