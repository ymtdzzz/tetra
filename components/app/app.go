package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ymtdzzz/tetra/adapter"
	"github.com/ymtdzzz/tetra/components/editor"
	"github.com/ymtdzzz/tetra/components/tree"
	"github.com/ymtdzzz/tetra/config"
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

	return Model{
		styles:  defaultStyles(),
		keyMap:  defaultKeyMap(),
		dbConns: conns,
		tree:    tree.New(conns),
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
	sidebar := m.styles.sidebar.Render(m.tree.View())
	mainTop := m.styles.mainTop.Render(m.editor.View())
	mainBottom := m.styles.mainBottom.Render("Main Bottom Panel")
	main := lipgloss.JoinVertical(lipgloss.Left, mainTop, mainBottom)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	return layout
}
