package main

import (
	"github.com/charmbracelet/lipgloss"
)

// TODO: ADD FUNCTIONALITY HERE
func (m *Model) MenuSelectDown() {
	if m.mainMenuFocus < Quit {
		m.mainMenuFocus++
	}
}

func (m *Model) MenuSelectUp() {
	if m.mainMenuFocus > Searching { // Don't go below Searching (0)
		m.mainMenuFocus--
	}
}

// TODO: FIX THE POSITION OF MainMenuView
func (m Model) MainMenuView() string {
	centerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center)

	return centerStyle.Render(m.list.View())
}
