package main

import (
	// "log"
	// "fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

var (

	// General.
	normal    = lipgloss.Color("#EEEEEE")
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	base      = lipgloss.NewStyle().Foreground(normal)

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render
	
	buttonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		MarginTop(1)

	activeButtonStyle = buttonStyle.
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		MarginRight(2).
		Underline(true)
)

func ConfigurationTitle(builder *strings.Builder, width int) *strings.Builder {

	descStyle := base.MarginTop(2).
		Align(lipgloss.Center).
		Width(width)

	infoStyle := base.
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(subtle).
		Width(width).
		Align(lipgloss.Center)

	title := figure.NewFigure("CONFIGURATION", "", true)

	desc := lipgloss.JoinVertical(lipgloss.Left,
		descStyle.Render(title.String()),
		infoStyle.Render("Customize Color Scheme, Editor Loading and other Defaults"),
	)

	row := lipgloss.JoinHorizontal(lipgloss.Top, desc)
	builder.WriteString(row + "\n\n")

	return builder
}

func ConfigurationOptions(builder *strings.Builder, width int) *strings.Builder {

	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		MarginTop(1)

	activeButtonStyle := buttonStyle.
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		MarginRight(2).
		Underline(true)

	mainMenu := activeButtonStyle.Render("MainMenu")
	searching := buttonStyle.Render("Searching")

	question := lipgloss.NewStyle().Width(width).Align(lipgloss.Left).Render("Startup Options: ")
	buttons := lipgloss.JoinHorizontal(lipgloss.Center, mainMenu, searching)

	ui := lipgloss.JoinVertical(lipgloss.Right, question, buttons)

	builder.WriteString(ui + "\n\n")
	return builder
}



/*
	TODO: We now have the ability to add buttons to the template and give them highlighting power
		Task 1: Build a Row of Buttons with interactive scrolling, i.e. h == left l == right 
		Task 2: Build 2 Rows of Buttons with interactive scrolling --> a [][]string should be used
		Task 3: Build the rest of the Button options and input functionality
*/


func (m *Model) CreateConfigurationTemplate(width, height int) string {
	doc := strings.Builder{}
	ConfigurationTitle(&doc, width)


	var style lipgloss.Style
	var buttons []string

	for _, opt := range m.OptionsMenu {
		style = activeButtonStyle
		buttons = append(buttons, style.Render(opt.buttonChoice))
	}
	
	row := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)

	return lipgloss.JoinVertical(lipgloss.Center, doc.String(), row)

}
