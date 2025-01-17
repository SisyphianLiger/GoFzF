package main

import ()

type ConfigurationChoices struct {
	width              int
	height             int
	ConfigurationFocus int
	OptionsMenu        []ButtonOptions
}

type ButtonOptions struct {
	buttonChoice string
	choices []string
}

// TODO: FIGURE OUT HOW INPUT WORKS...
type CustomInput struct {
	desc     string
	inputBox string
}

func (b ButtonOptions) FilterValue() string {
	return b.buttonChoice
}

func ConfigurationMenuData() []ButtonOptions {
	return []ButtonOptions{
		{buttonChoice: "Choose Startup Option: ", choices: []string{"Start Search", "MainMenu"}},
		{buttonChoice: "MainMenu"},
	}
}
