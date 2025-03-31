package result

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	showContextMenu key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		showContextMenu: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "show context menu"),
		),
	}
}
