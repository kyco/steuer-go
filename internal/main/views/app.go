package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"tax-calculator/internal/tax/views"
)

func Start() error {
	p := tea.NewProgram(
		views.NewRetroApp(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}