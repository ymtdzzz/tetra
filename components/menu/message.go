package menu

// FIXME: This message is handled in the app component.
// It should be moved to the menu component but import cycles prevent that.
type ShowMenuMsg struct {
	Items []Item
}
