package notification

import "github.com/charmbracelet/lipgloss"

type styles struct {
	notificationBase lipgloss.Style
}

func defaultStyles() styles {
	notificationBase := lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Border(lipgloss.ThickBorder()).
		BorderTop(false).
		BorderRight(false).
		BorderBottom(false).
		BorderLeft(true)

	return styles{
		notificationBase: notificationBase,
	}
}
