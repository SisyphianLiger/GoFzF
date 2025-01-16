package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

/*
	Particular Branches available in the TUI
*/
const (
	Searching     status = iota // 0
	Configuration               // 1
	Quit                        // 2
	MainMenu                    // 3
)

func main() {
	m := New(true)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
