package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Start() error {
	p := tea.NewProgram(
		NewAppModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}