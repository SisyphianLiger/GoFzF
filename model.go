package main

import (
	"github.com/charmbracelet/bubbles/list"
)

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

func (m *Model) initializeList(width, height int, fn Data, startingPoint status, title string) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = title
	m.list.SetShowStatusBar(false)
	m.list.SetItems(fn())
	m.mainMenuFocus = startingPoint
}
