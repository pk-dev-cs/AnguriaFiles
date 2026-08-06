package main

import (
	"charm.land/bubbles/v2/filepicker"
	"charm.land/lipgloss/v2"
)

func newPicker(directory string) filepicker.Model {
	picker := filepicker.New()

	picker.AutoHeight = false
	picker.CurrentDirectory = directory
	picker.ShowPermissions = true
	picker.ShowSize = true
	picker.ShowHidden = false
	picker.DirAllowed = true
	picker.FileAllowed = true
	
        picker.KeyMap.Open.SetKeys("enter")
        picker.KeyMap.Select.SetKeys("space")
	return picker
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
