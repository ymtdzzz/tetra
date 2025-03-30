package result

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	table table.Model
}

func New() Model {
	table := table.New([]table.Column{}).WithMaxTotalWidth(30).WithPageSize(5)

	return Model{
		table: table,
	}
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
			rows[i] = table.NewRow(r)
		}
		m.table = m.table.WithColumns(cols).WithRows(rows)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
	// return "No result"
}

func (m *Model) UpdateLayout(width, height int) {
}
