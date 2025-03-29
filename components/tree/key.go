package tree

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	focusSearch    key.Binding
	scrollToTop    key.Binding
	scrollToBottom key.Binding
	halfPageUp     key.Binding
	halfPageDown   key.Binding
	down           key.Binding
	up             key.Binding
	enter          key.Binding
	expand         key.Binding
	shrink         key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		focusSearch: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		scrollToTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "scroll to top"),
		),
		scrollToBottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "scroll to bottom"),
		),
		halfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp("u/^u", "½ page up"),
		),
		halfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp("d/^d", "½ page down"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down/j", "down"),
		),
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up/k", "up"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "expand/collapse"),
		),
		expand: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("right/l", "expand"),
		),
		shrink: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("left/h", "shrink"),
		),
	}
}
