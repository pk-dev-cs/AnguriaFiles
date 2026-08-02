package main

import (
	tea "charm.land/bubbletea/v2"
)

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
