package menu

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ymtdzzz/tetra/components"
)

type Item struct {
	Label    string
	Key      key.Binding
	Callback tea.Cmd
}

func (i Item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Label)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	focus bool
	items []Item
	list  list.Model
}

func New() Model {
	l := list.New([]list.Item{}, itemDelegate{}, 30, 10)
	l.DisableQuitKeybindings()
	l.SetFilteringEnabled(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)

	return Model{
		items: []Item{},
		list:  l,
	}
}

func (m *Model) Focus(focus bool) {
	m.focus = focus
}

func (m Model) Focused() bool {
	return m.focus
}

func (m *Model) SetItems(items []Item) {
	lis := make([]list.Item, len(items))
	for i, item := range items {
		lis[i] = item
	}
	m.list.SetItems(lis)
	m.items = items
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for _, item := range m.items {
		if k, ok := msg.(tea.KeyMsg); ok && key.Matches(k, item.Key) {
			cmds = append(cmds, tea.Batch(
				item.Callback,
				func() tea.Msg {
					return components.CloseMenuMsg{}
				},
			))
			break
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			cmds = append(cmds, tea.Batch(
				m.items[m.list.Index()].Callback,
				func() tea.Msg {
					return components.CloseMenuMsg{}
				},
			))
		} else if msg.String() == "esc" {
			cmds = append(cmds, func() tea.Msg {
				return components.CloseMenuMsg{}
			})
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.focus {
		return ""
	}
	return m.list.View()
}

func (m Model) UpdateLayout(width, height int) {
	m.list.SetWidth(width)
	m.list.SetHeight(height)
}
