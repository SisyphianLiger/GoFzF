package main

import (

	tea "github.com/charmbracelet/bubbletea"
)


type StateResult struct {
	status status
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	m.startUp(msg)

	switch m.state {
	case MainMenu:
		m.MainMenuState(msg)

		if m.MainMenuBranch(msg).status == Quit {
			return m, tea.Quit	
		}

	case Searching:
	// TODO: Loads Searching Module

	case Configuration:
	// TODO: Put Configuration into View

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}


func (m *Model) startUp(msg tea.Msg) {

	if !m.sizeOfScreen {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// Assiging this for other options
			m.width, m.height = msg.Width, msg.Height
			m.sizeOfScreen = true
		}
	}
}


func (m *Model) MainMenuState(msg tea.Msg) {
	// TODO: INIT Main Menu
	if !m.intialState {
		m.initializeList(m.width, m.height, MainMenuOptions, Searching, "Main Menu")
		m.intialState = true
	}
}


func (m *Model) MainMenuBranch(msg tea.Msg) StateResult {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "j", "down":
			m.MenuSelectDown()

		case "k", "up":
			m.MenuSelectUp()

		case "q":
			m.state = Quit
			return StateResult{ status: Quit }

		case "enter":
			if m.mainMenuFocus == Searching {
				m.state = Searching
			}

			if m.mainMenuFocus == Configuration {
				m.state = Configuration
				m.initializeList(m.width, m.height, ConfigurationMenuData, Configuration, "Configuration")
			}

			if m.mainMenuFocus == Quit {
				m.state = Quit
				return StateResult{ status: Quit }
			}
		}
	}

	return StateResult{}
}
