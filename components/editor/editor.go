package editor

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ymtdzzz/tetra/adapter"
	"github.com/ymtdzzz/tetra/components/result"
)

type Model struct {
	keyMap   keyMap
	textarea textarea.Model
	focus    bool
	conn     *adapter.DBConnection
}

func New() Model {
	ta := textarea.New()

	return Model{
		keyMap:   defaultKeyMap(),
		textarea: ta,
	}
}

func (m *Model) Focus(focus bool) {
	m.focus = focus
	if focus {
		m.textarea.Focus()
	} else {
		m.textarea.Blur()
	}
}

func (m Model) Focused() bool {
	return m.focus
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.textarea.Focus(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case SetConnMsg:
		m.conn = msg.Conn
	}

	if !m.focus {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.executeQuery):
			if m.conn == nil {
				break
			}
			cmds = append(cmds, func() tea.Msg {
				// TODO: context cancel
				res, err := m.conn.Adapter.RunQuery(context.Background(), m.textarea.Value())
				// TODO: handle error
				if err != nil {
					panic(err)
				}
				if r, ok := res.([]map[string]any); ok {
					return result.QueryResultMsg{
						Result: r,
					}
				}
				return nil
			})
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.conn == nil {
		style := lipgloss.NewStyle().
			Width(m.textarea.Width() + 6).
			Height(m.textarea.Height())
		return style.Render("No Database connection selected\nPress '^e' in the tree to select a connection")
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		fmt.Sprintf("Connected to %s", m.conn.Name),
		m.textarea.View(),
	)
}

func (m *Model) UpdateLayout(width, height int) {
	m.textarea.SetWidth(width)
	m.textarea.SetHeight(height - 1)
}
