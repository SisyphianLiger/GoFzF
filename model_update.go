package main

import (
	// "log"

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
		if m.ConfigurationBranch(msg).status == Quit {
			return m, tea.Quit
		}
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
			return StateResult{status: Quit}

		case "enter":
			if m.mainMenuFocus == Searching {
				m.state = Searching
			}

			if m.mainMenuFocus == Configuration {
				// Change State, Add Data
				m.state = Configuration
				m.ConfigurationChoices.OptionsMenu = ConfigurationMenuData()

			}

			if m.mainMenuFocus == Quit {
				m.state = Quit
				return StateResult{status: Quit}
			}
		}
	}

	return StateResult{}
}

func (m *Model) ConfigurationBranch(msg tea.Msg) StateResult {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return StateResult{status: Quit}
		case "right", "l", "n", "tab":
			m.ConfigurationFocus = min(m.ConfigurationFocus+1, len(m.OptionsMenu)-1)
		case "left", "h", "p", "shift+tab":
			m.ConfigurationFocus = max(m.ConfigurationFocus-1, 0)
		// Here Again like with MainMenu we are able to go up and down and enter does stuff
		case "enter":
			if m.ConfigurationFocus == 0 {
				m.state = MainMenu
				return StateResult{status: MainMenu}	
			}
		}
	}

	return StateResult{}
}
