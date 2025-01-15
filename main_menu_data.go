package main

import (
	"github.com/charmbracelet/bubbles/list"
)

type status int

type MenuOption struct {
	number      int
	title       string
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

func MainMenuOptions() []list.Item {
	return []list.Item{
		MenuOption{title: "Start Searching"},
		MenuOption{title: "Configuration Options"},
		MenuOption{title: "Quit"},
	}
}
