package main

import (
	"github.com/charmbracelet/bubbles/list"
)

type ConfigurationChoices struct {
	width              int
	height             int
	ConfigurationFocus status
}

type Toggle struct {
	on        bool
	choiceOne string
	choiceTwo string
}

// TODO: FIGURE OUT HOW INPUT WORKS...
type CustomInput struct {
	desc     string
	inputBox string
}

type ButtonOptions struct {
	buttonName string
	Toggle
	CustomInput
}

func (b ButtonOptions) FilterValue() string {
	return b.buttonName
}

func ConfigurationMenuData() []list.Item {
	return []list.Item{
		ButtonOptions{buttonName: "MainMenu"},
	}
}
