package main

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
   This will be the main Menu. Will Be uploaded when for the Users Fist Experience
   After the first run the User should should be able to access it from the help menu with H
*/

const (
	Searching     status = iota // 0
	Configuration               // 1
	Quit                        // 2
	MainMenu                    // 3
)

type ConfigurationChoices struct {
	width              int
	height             int
	ConfigurationFocus status
}

type Model struct {
	list          list.Model
	err           error
	sizeOfScreen  bool   // Used To Grab Window Component so the options are displayed properly
	state         status // Determines Branches of Options / Where Initializer begins
	mainMenuFocus status //
	ConfigurationChoices
	intialState bool // Used to Determine what case we load for the model on initializtion
}

func New(initialState bool) *Model {
	if initialState {
		return &Model{
			state: MainMenu,
		}
	}
	return &Model{
		state: Searching,
	}
}

// TODO: ADD FUNCTIONALITY HERE
func (m *Model) SelectDown() {
	if m.mainMenuFocus < Quit {
		m.mainMenuFocus++
	}
}

func (m *Model) SelectUp() {
	if m.mainMenuFocus > Searching { // Don't go below Searching (0)
		m.mainMenuFocus--
	}
}

type Data func() []list.Item

/*
Generalized Initial List Creator
*/
func (m *Model) initializeList(width, height int, fn Data, startingPoint status, title string) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = title
	m.list.SetShowStatusBar(false)
	m.list.SetItems(fn())
	m.mainMenuFocus = startingPoint
}

func (m Model) Init() tea.Cmd {
	return nil
}

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
				m.SelectDown()
			case "k", "up":
				m.SelectUp()
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
					m.initializeList(m.width, m.height, ConfigurationMenuData, Configuration, "Configuration????")
				}
				if m.mainMenuFocus == Quit {
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

func (m Model) View() string {
	// When we exit the program we make sure to render nothing to the screen
	// therefore it will return to normal
	if m.state == Searching {
		return "Start Search!"
	}
	if m.state == Configuration {
		return CreateConfigurationTemplate(m.width, m.height)
	}
	if m.state == Quit {
		return ""
	}

	return m.MainMenuView()
}

// TODO: FIX THE POSITION OF MainMenuView
func (m Model) MainMenuView() string {
	centerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)

	return centerStyle.Render(m.list.View())
}

func main() {
	m := New(true)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
