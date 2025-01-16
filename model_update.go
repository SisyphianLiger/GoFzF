package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if !m.sizeOfScreen {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// Assiging this for other options
			m.width, m.height = msg.Width, msg.Height
			m.sizeOfScreen = true
		}
	}

	switch m.state {
	case MainMenu:
		// TODO: INIT Main Menu
		if !m.intialState {
			m.initializeList(m.width, m.height, MainMenuOptions, Searching, "Main Menu")
			m.intialState = true
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j", "down":
				m.MenuSelectDown()
			case "k", "up":
				m.MenuSelectUp()
			case "q":
				m.state = Quit
				return m, tea.Quit

			case "enter":
				// TODO: Implement screen differences here
				if m.mainMenuFocus == Searching {
					m.state = Searching
					// Logic to Search Here
				}
				if m.mainMenuFocus == Configuration {
					m.state = Configuration
					m.initializeList(m.width, m.height, ConfigurationMenuData, Configuration, "Configuration")
				}
				if m.mainMenuFocus == Quit {
					m.state = Quit
					return m, tea.Quit
				}
			}
		}

	case Searching:
	// TODO: Loads Searching Module

	case Configuration:
	// TODO: Put Configuration into View

	case Quit:

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
