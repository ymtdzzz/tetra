package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ymtdzzz/tetra/adapter"
	"github.com/ymtdzzz/tetra/components/tree"
	"github.com/ymtdzzz/tetra/config"
)

type Model struct {
	keyMap  keyMap
	dbConns adapter.DBConnections
	tree    tree.Model
}

func New() (Model, error) {
	config, err := config.LoadConfig("./config.toml")
	if err != nil {
		return Model{}, err
	}
	conns := adapter.NewDBConnections(config)

	return Model{
		keyMap:  defaultKeyMap(),
		dbConns: conns,
		tree:    tree.New(conns),
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
	}
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.tree.View()
	// return fmt.Sprint(m.msg)
}
