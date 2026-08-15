package main

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"path/filepath"

	"charm.land/bubbles/v2/filepicker"
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

func copyFile(sourcePath, destinationPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dialog != nil {
		switch typedMsg := msg.(type) {
		case tea.KeyPressMsg:
			switch typedMsg.String() {
			case "enter", "y":
				return m, m.closeDialog(true)
			case "esc", "n":
				return m, m.closeDialog(false)
			}

			return m, nil
		case tea.WindowSizeMsg:
			m.width = typedMsg.Width
			m.height = typedMsg.Height
			m.resizePickers()
			return m, nil
		}
	}

	if m.popupMessage != "" {
		switch typedMsg := msg.(type) {
		case tea.KeyPressMsg:
			if typedMsg.String() == "esc" {
				m.popupMessage = ""
			}

			return m, nil

		case tea.WindowSizeMsg:
			m.width = typedMsg.Width
			m.height = typedMsg.Height
			m.resizePickers()
		}
	}

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
		case "delete":
			picker := m.activePicker()
			path := picker.HighlightedPath()
			if path != "" {
				m.showDialog(
				"Potwierdzenie",
				"Czy chcesz kontynuować?",
				func(currentModel *model, result bool) tea.Cmd {
					if result {
						err := os.Remove(path)
						if err != nil{
							currentModel.showPopup(err.Error())							
						}
					}

					return picker.Init()
				},
			) }

			return m, nil	

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
