package main

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

/*
   This will be the main Menu. Will Be uploaded when for the Users Fist Experience
   After the first run the User should should be able to access it from the help menu with H
*/
type status int

const (

)
type MenuOption struct {
	number int
	title string
	description string 
}

func (menu MenuOption) FilterValue() string {
	return menu.title
}

func (menu MenuOption) Title() string {
	return menu.title
}

func (menu MenuOption) Description() string {
	return menu.description
}


type Model struct {
	list list.Model
	err error
	focus status
	quit bool 
}

func New() *Model{
	return &Model{}
}

// TODO: ADD FUNCTIONALITY HERE
func (m *Model) SelectDown()  {
	if m.focus == 3 {
		m.focus = status(3)
		return
	}
	m.focus++
	return
}

func (m *Model) SelectUp()  {
	if m.focus == 0 {
		m.focus = status(0)
		return
	}
	m.focus--
	return
}

type Data func() []list.Item

func (m *Model) initMenuList(width, height int, fn Data) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "Main Menu"
	m.list.SetItems(fn())
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initMenuList(msg.Width, msg.Height, MainMenu)
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.SelectDown()
		case "k", "up":
			m.SelectUp()
		case "enter":

			if m.focus == 0 {
			
			}
			if m.focus == 1 {
			
			}
			if m.focus == 2 {
			
			}
			if m.focus == 3 {
				m.quit = true
				return m, tea.Quit
			}
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// When we exit the program we make sure to render nothing to the screen 
	// therefore it will return to normal
	if m.quit == true {
		return ""
	}

	return m.list.View()
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _,err := p.Run(); err != nil {
		os.Exit(1)
	}
}
