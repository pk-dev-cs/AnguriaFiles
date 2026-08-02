package main

import (
	"strings"

	"charm.land/lipgloss/v2"
)

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
