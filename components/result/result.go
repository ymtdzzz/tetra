package result

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/ymtdzzz/tetra/adapter"
	"github.com/ymtdzzz/tetra/components/menu"
	"github.com/ymtdzzz/tetra/components/notification"
)

type Model struct {
	keyMap        keyMap
	focus         bool
	table         table.Model
	rows          []table.Row
	width, height int
}

func New() Model {
	t := table.New([]table.Column{}).WithMaxTotalWidth(30).WithPageSize(5)

	return Model{
		keyMap: defaultKeyMap(),
		table:  t,
		rows:   []table.Row{},
	}
}

func (m *Model) Focus(focus bool) {
	m.focus = focus
	m.table = m.table.Focused(focus)
}

func (m Model) Focused() bool {
	return m.focus
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case QueryResultMsg:
		if len(msg.Result) == 0 {
			break
		}
		cols := []table.Column{}
		rows := make([]table.Row, len(msg.Result))
		for k := range msg.Result[0] {
			cols = append(cols, table.NewColumn(k, k, 10))
		}
		for i, r := range msg.Result {
			for k := range r {
				if v, ok := r[k].([]byte); ok {
					r[k] = string(v)
				}
			}
			rows[i] = table.NewRow(r)
		}
		m.table = m.table.WithColumns(cols).WithRows(rows)
		m.rows = rows
	}

	if !m.focus {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.showContextMenu):
			cmds = append(cmds, m.showContextMenuCmd())
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.hasResult() {
		style := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height)
		return style.Render("No results")
	}

	return m.table.View()
}

func (m *Model) UpdateLayout(width, height int) {
	m.height = height
	m.width = width

	m.table = m.table.WithMaxTotalWidth(width)
	m.table = m.table.WithPageSize(height - 6)
}

func (m Model) showContextMenuCmd() tea.Cmd {
	return func() tea.Msg {
		items := []menu.Item{}

		if m.hasResult() {
			items = append(items, menu.Item{
				Label: "Copy as CSV",
				Key:   key.NewBinding(key.WithKeys("y")),
				Callback: func() tea.Msg {
					err := clipboard.WriteAll(adapter.ConvertResultToCSV(m.rows))
					msg := "Copied as CSV!"
					if err != nil {
						msg = fmt.Sprintf("Failed to copy as CSV: %s", err)
					}
					return notification.NotificationMsg{
						Message: msg,
					}
				},
			})
		}

		return menu.ShowMenuMsg{
			Items: items,
		}
	}
}

func (m Model) hasResult() bool {
	return m.table.TotalRows() > 0
}
