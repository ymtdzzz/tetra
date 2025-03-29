package app

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}
