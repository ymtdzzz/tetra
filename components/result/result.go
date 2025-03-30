package result

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	focus         bool
	table         table.Model
	width, height int
}

func New() Model {
	table := table.New([]table.Column{}).WithMaxTotalWidth(30).WithPageSize(5)

	return Model{
		table: table,
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
	}

	if !m.focus {
		return m, tea.Batch(cmds...)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.table.TotalRows() == 0 {
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
