package main

import (
	"os"
	"os/exec"
	"runtime"
	tea "charm.land/bubbletea/v2"
)

type fileOpenedMsg struct {
	path string
	err  error
}

func openFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command(
				"rundll32",
				"url.dll,FileProtocolHandler",
				path,
			)

		case "darwin":
			cmd = exec.Command("open", path)

		default:
			cmd = exec.Command("xdg-open", path)
		}

		return fileOpenedMsg{
			path: path,
			err:  cmd.Run(),
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typedMsg.Width
		m.height = typedMsg.Height
		m.resizePickers()

	case fileOpenedMsg:
		if typedMsg.err != nil {
			m.error = typedMsg.err
		}

		return m, nil

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

		case "space":
			picker := m.activePicker()
			path := picker.HighlightedPath()
			if path == "" {
				return m, nil
			}

			picker.Path = path
			m.selectedFile = path
			m.resizePickers()

			return m, nil

		case "enter":
			picker := m.activePicker()
			path := picker.HighlightedPath()
			if path == "" {
				return m, nil
			}

			info, err := os.Stat(path)
			if err != nil {
				return m, nil
			}

			if info.IsDir() {
				updatedPicker, cmd := picker.Update(msg)
				*picker = updatedPicker

				return m, cmd
			}

			return m, openFileCmd(path)
		}

		picker := m.activePicker()

		updatedPicker, cmd := picker.Update(msg)
		*picker = updatedPicker

		return m, cmd
	}

	var leftCmd tea.Cmd
	var rightCmd tea.Cmd

	m.leftPicker, leftCmd = m.leftPicker.Update(msg)
	m.rightPicker, rightCmd = m.rightPicker.Update(msg)

	return m, tea.Batch(leftCmd, rightCmd)
}
