package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/filepicker"
	tea "charm.land/bubbletea/v2"
)

type pane int

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

	popupMessage string
	dialog       *dialogState
	error        error
}

type dialogState struct {
	title     string
	message   string
	onResolve func(*model, bool) tea.Cmd
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

func (m *model) showPopup(message string) {
	m.popupMessage = message
}

func (m *model) showDialog(
	title, message string,
	onResolve func(*model, bool) tea.Cmd,
) {
	m.dialog = &dialogState{
		title:     title,
		message:   message,
		onResolve: onResolve,
	}
}

func (m *model) closeDialog(result bool) tea.Cmd {
	if m.dialog == nil {
		return nil
	}

	onResolve := m.dialog.onResolve
	m.dialog = nil

	if onResolve != nil {
		return onResolve(m, result)
	}

	return nil
}
