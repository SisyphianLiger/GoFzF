package main

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (

	// General.
	normal    = lipgloss.Color("#EEEEEE")
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	base = lipgloss.NewStyle().Foreground(normal)

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

	// Title Styling
	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	descStyle = base.MarginTop(1).
			Align(lipgloss.Center)

	infoStyle = base.
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)
)

func ConfigurationTitle(builder *strings.Builder) *strings.Builder {
	var title strings.Builder
	desc := lipgloss.JoinVertical(lipgloss.Left,
		descStyle.Render("Configuration Options:"),
		infoStyle.Render("Customize Color Scheme, Editor Loading and other Defaults"),
	)

	row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
	builder.WriteString(row + "\n\n")

	return builder
}

/*
Step One: Create the document that will hold the "strings"
*/

func CreateConfigurationTemplate(width, height int) string {
	doc := strings.Builder{}

	ConfigurationTitle(&doc)
	docStyle := lipgloss.NewStyle().Padding(1, 2, 1, 2)

	return docStyle.Render(doc.String())

}
