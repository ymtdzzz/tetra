package app

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit          key.Binding
	focusToDBNavi key.Binding
	focusToEditor key.Binding
	focusToResult key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("^c", "quit"),
		),
		focusToDBNavi: key.NewBinding(
			key.WithKeys("1"),
		),
		focusToEditor: key.NewBinding(
			key.WithKeys("2"),
		),
		focusToResult: key.NewBinding(
			key.WithKeys("3"),
		),
	}
}
