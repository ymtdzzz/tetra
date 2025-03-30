package editor

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	executeQuery key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		executeQuery: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("^j", "Execute query"),
		),
	}
}
