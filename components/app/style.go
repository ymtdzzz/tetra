package app

import "github.com/charmbracelet/lipgloss"

type styles struct {
	sidebar    lipgloss.Style
	mainTop    lipgloss.Style
	mainBottom lipgloss.Style
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

	return styles{
		sidebar:    sidebar,
		mainTop:    mainTop,
		mainBottom: mainBottom,
	}
}
