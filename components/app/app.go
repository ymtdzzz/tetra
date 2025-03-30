package app

import (
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/ymtdzzz/tetra/adapter"
	"github.com/ymtdzzz/tetra/components/editor"
	"github.com/ymtdzzz/tetra/components/tree"
	"github.com/ymtdzzz/tetra/config"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	styles  styles
	keyMap  keyMap
	dbConns adapter.DBConnections
	tree    tree.Model
	editor  editor.Model
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
	}, nil
}

func (m Model) Close() error {
	return m.dbConns.Close()
}

func (m Model) Init() tea.Cmd {
	return m.tree.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		height, width := msg.Height, msg.Width
		sidebarWidth := int(float64(width) * 0.3)
		mainWidth := width - sidebarWidth
		mainTopHeight := height / 2
		mainBottomHeight := height - mainTopHeight

		m.styles.mainTop = m.styles.mainTop.
			Width(mainWidth - m.styles.mainTop.GetHorizontalFrameSize()).
			Height(mainTopHeight - m.styles.mainTop.GetVerticalFrameSize())
		m.styles.mainBottom = m.styles.mainBottom.
			Width(mainWidth - m.styles.mainBottom.GetHorizontalFrameSize()).
			Height(mainBottomHeight - m.styles.mainBottom.GetVerticalFrameSize())

		m.tree.UpdateLayout(
			sidebarWidth-m.styles.sidebar.GetHorizontalFrameSize(),
			height-m.styles.sidebar.GetVerticalFrameSize(),
		)
	}
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	m.styles.sidebar = focusedStyle(m.styles.sidebar, m.tree.Focused())

	sidebar := m.styles.sidebar.Render(m.tree.View())
	sidebar = renderWithTitle(sidebar, "DB Navigator [1]", m.styles.sidebar)
	mainTop := m.styles.mainTop.Render(m.editor.View())
	mainBottom := m.styles.mainBottom.Render("Main Bottom Panel")
	main := lipgloss.JoinVertical(lipgloss.Left, mainTop, mainBottom)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	return layout
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
