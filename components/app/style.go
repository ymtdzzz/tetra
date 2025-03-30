package app

import "github.com/charmbracelet/lipgloss"

type styles struct {
	sidebar    lipgloss.Style
	mainTop    lipgloss.Style
	mainBottom lipgloss.Style
	menu       lipgloss.Style
}

func defaultStyles() styles {
	sidebar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	mainTop := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	mainBottom := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	menu := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	return styles{
		sidebar:    sidebar,
		mainTop:    mainTop,
		mainBottom: mainBottom,
		menu:       menu,
	}
}

func focusedStyle(style lipgloss.Style, focused bool) lipgloss.Style {
	if !focused {
		return style
	}
	newStyle := style
	return newStyle.BorderForeground(lipgloss.Color("#FFCC00"))
}
