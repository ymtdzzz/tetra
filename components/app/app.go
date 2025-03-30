package app

import (
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/ymtdzzz/tetra/adapter"
	"github.com/ymtdzzz/tetra/components"
	"github.com/ymtdzzz/tetra/components/editor"
	"github.com/ymtdzzz/tetra/components/result"
	"github.com/ymtdzzz/tetra/components/tree"
	"github.com/ymtdzzz/tetra/config"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const (
	PANE_DB_NAVIGATOR = iota
	PANE_EDITOR
)

type Model struct {
	styles  styles
	keyMap  keyMap
	dbConns adapter.DBConnections
	tree    tree.Model
	editor  editor.Model
	result  result.Model
}

func New() (Model, error) {
	config, err := config.LoadConfig("./config.toml")
	if err != nil {
		return Model{}, err
	}
	conns := adapter.NewDBConnections(config)

	tree := tree.New(conns)
	tree.Focus(true)

	return Model{
		styles:  defaultStyles(),
		keyMap:  defaultKeyMap(),
		dbConns: conns,
		tree:    tree,
		editor:  editor.New(),
		result:  result.New(),
	}, nil
}

func (m Model) Close() error {
	return m.dbConns.Close()
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.tree.Init(),
		m.editor.Init(),
		m.result.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.focusToDBNavi):
			m.focusPane(PANE_DB_NAVIGATOR)
			return m, nil
		case key.Matches(msg, m.keyMap.focusToEditor):
			m.focusPane(PANE_EDITOR)
			return m, nil
		}
	case tea.WindowSizeMsg:
		height, width := msg.Height, msg.Width
		sidebarWidth := int(float64(width) * 0.3)
		mainWidth := width - sidebarWidth
		mainTopHeight := height / 2
		mainBottomHeight := height - mainTopHeight

		m.tree.UpdateLayout(
			sidebarWidth-m.styles.sidebar.GetHorizontalFrameSize(),
			height-m.styles.sidebar.GetVerticalFrameSize(),
		)
		m.editor.UpdateLayout(
			mainWidth-m.styles.mainTop.GetHorizontalFrameSize(),
			mainTopHeight-m.styles.mainTop.GetVerticalFrameSize(),
		)
		m.result.UpdateLayout(
			mainWidth-m.styles.mainBottom.GetHorizontalFrameSize(),
			mainBottomHeight-m.styles.mainBottom.GetVerticalFrameSize(),
		)
	case components.FocusPaneEditorMsg:
		m.focusPane(PANE_EDITOR)
	}
	m.tree, cmd = m.tree.Update(msg)
	cmds = append(cmds, cmd)
	m.editor, cmd = m.editor.Update(msg)
	cmds = append(cmds, cmd)
	m.result, cmd = m.result.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.styles.sidebar = focusedStyle(m.styles.sidebar, m.tree.Focused())
	m.styles.mainTop = focusedStyle(m.styles.mainTop, m.editor.Focused())

	sidebar := m.styles.sidebar.Render(m.tree.View())
	sidebar = renderWithTitle(sidebar, "DB Navigator [1]", m.styles.sidebar)
	mainTop := m.styles.mainTop.Render(m.editor.View())
	mainTop = renderWithTitle(mainTop, "Editor [2]", m.styles.mainTop)
	mainBottom := m.styles.mainBottom.Render(m.result.View())
	main := lipgloss.JoinVertical(lipgloss.Left, mainTop, mainBottom)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	return layout
}

func (m *Model) focusPane(pane int) {
	switch pane {
	case PANE_DB_NAVIGATOR:
		m.tree.Focus(true)
		m.editor.Focus(false)
	case PANE_EDITOR:
		m.tree.Focus(false)
		m.editor.Focus(true)
	}
}

func renderWithTitle(view, title string, style lipgloss.Style) string {
	lines := strings.Split(view, "\n")
	if len(lines) == 0 {
		return view
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(style.GetBorderTopForeground())

	plain := stripansi.Strip(lines[0])
	var b strings.Builder
	titleWidth := runewidth.StringWidth(title) + 2
	for i := 0; i < titleWidth; i++ {
		b.WriteString(style.GetBorderStyle().Top)
	}
	replaced := strings.Replace(plain, b.String(), " "+title+" ", 1)

	lines[0] = titleStyle.Render(replaced)

	return strings.Join(lines, "\n")
}
