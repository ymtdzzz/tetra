package tree

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ymtdzzz/tetra/adapter"
	"github.com/ymtdzzz/tetra/components"
	"github.com/ymtdzzz/tetra/components/editor"
)

const (
	NODE_TYPE_CONNECTION = iota
	NODE_TYPE_DATABASE_FOLDER
	NODE_TYPE_DATABASE
	NODE_TYPE_TABLE_FOLDER
	NODE_TYPE_TABLE

	LABEL_DUMMY = "this_is_dummy_label"
)

type TreeNode struct {
	Icon     rune
	Label    string
	Children []*TreeNode
	Expanded bool
	Parent   *TreeNode
	nodeType int
	loading  bool
	conn     *adapter.DBConnection
}

type Model struct {
	keyMap       keyMap
	Roots        []*TreeNode
	FlattenNodes []*TreeNode
	Cursor       int
	conns        adapter.DBConnections
	spinner      spinner.Model
	viewport     viewport.Model
	ready        bool
	textInput    textinput.Model
	query        string
	focus        bool
}

func New(conns adapter.DBConnections) Model {
	roots := make([]*TreeNode, len(conns))
	for i, conn := range conns {
		roots[i] = &TreeNode{
			Icon:     components.ICON_MYSQL,
			Label:    fmt.Sprintf("%s %s:%d", conn.Name, conn.DBConfig.Host, conn.DBConfig.Port),
			nodeType: NODE_TYPE_CONNECTION,
			conn:     conn,
		}
		roots[i].Children = []*TreeNode{
			{
				Label:  LABEL_DUMMY,
				Parent: roots[i],
			},
		}
	}

	ti := textinput.New()
	ti.Placeholder = "/ to search"
	ti.Prompt = "Query: "

	m := Model{
		keyMap:    defaultKeyMap(),
		Roots:     roots,
		conns:     conns,
		spinner:   spinner.New(),
		textInput: ti,
	}
	m.updateView(true)
	return m
}

func (m *Model) Focus(focus bool) {
	m.focus = focus
}

func (m Model) Focused() bool {
	return m.focus
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		textinput.Blink,
	)
}

func flattenAll(roots []*TreeNode, query string) []*TreeNode {
	nodes := []*TreeNode{}
	for _, r := range roots {
		flatten(r, &nodes)
	}
	if query == "" {
		return nodes
	}
	result := []*TreeNode{}
	for _, n := range nodes {
		if n.nodeType != NODE_TYPE_TABLE || strings.Contains(n.Label, query) {
			result = append(result, n)
		}
	}
	return result
}

func flatten(n *TreeNode, acc *[]*TreeNode) {
	*acc = append(*acc, n)
	if n.Expanded {
		for _, c := range n.Children {
			flatten(c, acc)
		}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case connectionOpenMsg:
		// TODO: Error handling
		cmds = append(cmds, func() tea.Msg {
			_ = msg.conn.Adapter.Open()
			return connectionOpenDoneMsg(msg)
		})
	case connectionOpenDoneMsg:
		node := msg.node.Parent
		fdatabase := &TreeNode{
			Icon:     components.ICON_FOLDER,
			Label:    "Databases",
			Parent:   node,
			nodeType: NODE_TYPE_DATABASE_FOLDER,
			conn:     msg.conn,
		}
		databases := &TreeNode{
			Label:  LABEL_DUMMY,
			Parent: fdatabase,
			conn:   msg.conn,
		}
		fdatabase.Children = []*TreeNode{databases}
		node.Children = []*TreeNode{
			fdatabase,
			// TODO: Users, System Information, etc.
		}
		m.updateView(true)
	case listDatabasesMsg:
		cmds = append(cmds, func() tea.Msg {
			// TODO: Error handling
			databases, _ := msg.conn.Adapter.ListDatabases()
			return listDatabasesDoneMsg{
				conn:      msg.conn,
				databases: databases,
				node:      msg.node,
			}
		})
	case listDatabasesDoneMsg:
		databases := msg.databases
		node := msg.node.Parent
		dbnodes := make([]*TreeNode, len(databases))
		for i, db := range databases {
			dbnode := &TreeNode{
				Icon:     components.ICON_DATABASE,
				Label:    db,
				Parent:   node,
				nodeType: NODE_TYPE_DATABASE,
				conn:     msg.conn,
			}
			ftable := &TreeNode{
				Icon:     components.ICON_FOLDER,
				Label:    "Tables",
				Parent:   dbnode,
				nodeType: NODE_TYPE_TABLE_FOLDER,
				conn:     msg.conn,
			}
			dbnode.Children = []*TreeNode{ftable}
			tables := &TreeNode{
				Label:  LABEL_DUMMY,
				Parent: ftable,
				conn:   msg.conn,
			}
			ftable.Children = []*TreeNode{tables}
			dbnodes[i] = dbnode
		}
		node.Children = dbnodes
		m.updateView(true)
	case listTablesMsg:
		cmds = append(cmds, func() tea.Msg {
			// TODO: Error handling
			tables, _ := msg.conn.Adapter.ListTables(msg.database)
			return listTablesDoneMsg{
				conn:   msg.conn,
				tables: tables,
				node:   msg.node,
			}
		})
	case listTablesDoneMsg:
		tables := msg.tables
		node := msg.node.Parent
		tableNodes := make([]*TreeNode, len(tables))
		for i, table := range tables {
			tableNodes[i] = &TreeNode{
				Icon:     components.ICON_TABLE,
				Label:    table,
				Parent:   node,
				nodeType: NODE_TYPE_TABLE,
				conn:     msg.conn,
			}
			// TODO: Columns, Constraints, Indexes, etc.
		}
		node.Children = tableNodes
		m.updateView(true)
	}

	if !m.focus {
		return m, tea.Batch(cmds...)
	}

	if m.textInput.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case msg.String() == "enter":
				m.query = m.textInput.Value()
				m.updateView(true)
				m.textInput.Blur()
			case msg.String() == "esc":
				m.textInput.Blur()
			}
		}

		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
		m.updateView(false)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.focusSearch):
			cmds = append(cmds, m.textInput.Focus())
		case key.Matches(msg, m.keyMap.scrollToTop):
			m.Cursor = 0
			m.viewport.GotoTop()
			m.updateView(false)
		case key.Matches(msg, m.keyMap.scrollToBottom):
			m.Cursor = len(m.FlattenNodes) - 1
			m.viewport.GotoBottom()
			m.updateView(false)
		case key.Matches(msg, m.keyMap.halfPageUp):
			m.Cursor -= m.viewport.Height / 2
			if m.Cursor < 0 {
				m.Cursor = 0
			}
			m.viewport.HalfPageUp()
			m.updateView(false)
		case key.Matches(msg, m.keyMap.halfPageDown):
			m.Cursor += m.viewport.Height / 2
			if m.Cursor >= len(m.FlattenNodes) {
				m.Cursor = len(m.FlattenNodes) - 1
			}
			m.viewport.HalfPageDown()
			m.updateView(false)
		case key.Matches(msg, m.keyMap.down):
			if m.Cursor < len(m.FlattenNodes)-1 {
				m.Cursor++
				if m.Cursor >= m.viewport.Height+m.viewport.YOffset-5 {
					m.viewport.SetYOffset(m.viewport.YOffset + 1)
				}
			}
			m.updateView(false)
		case key.Matches(msg, m.keyMap.scrollRight):
			m.viewport.ScrollRight(1)
			m.updateView(false)
		case key.Matches(msg, m.keyMap.scrollLeft):
			m.viewport.ScrollLeft(1)
			m.updateView(false)
		case key.Matches(msg, m.keyMap.up):
			if m.Cursor > 0 {
				m.Cursor--
				if m.Cursor < m.viewport.YOffset+1 {
					m.viewport.SetYOffset(m.viewport.YOffset - 1)
				}
			}
			m.updateView(false)
		case key.Matches(msg, m.keyMap.enter):
			node := m.FlattenNodes[m.Cursor]
			node.Expanded = !node.Expanded
			if node.Expanded {
				cmds = append(cmds, m.handleExpand(node))
			}
			m.updateView(true)
		case key.Matches(msg, m.keyMap.selectNode):
			node := m.FlattenNodes[m.Cursor]
			cmds = append(cmds, tea.Batch(
				func() tea.Msg {
					return editor.SetConnMsg{
						Conn: node.conn,
					}
				},
				func() tea.Msg {
					return components.FocusPaneEditorMsg{}
				},
			))
		}
	}

	return m, tea.Batch(cmds...)
}

func indent(n *TreeNode) int {
	level := 0
	for p := n.Parent; p != nil; p = p.Parent {
		level++
	}
	return level
}

func (m *Model) handleExpand(node *TreeNode) tea.Cmd {
	switch node.nodeType {
	case NODE_TYPE_CONNECTION:
		if !node.conn.Adapter.Status().Opened {
			node.Children[0].loading = true
			return func() tea.Msg {
				return connectionOpenMsg{
					conn: node.conn,
					node: node.Children[0],
				}
			}
		}
	case NODE_TYPE_DATABASE_FOLDER:
		if !node.conn.Adapter.Status().DatabaseLoaded {
			node.Children[0].loading = true
			return func() tea.Msg {
				return listDatabasesMsg{
					conn: node.conn,
					node: node.Children[0],
				}
			}
		}
	case NODE_TYPE_TABLE_FOLDER:
		ok, loaded := node.conn.Adapter.Status().TableLoaded[node.Parent.Label]
		if !ok || !loaded {
			node.Children[0].loading = true
			return func() tea.Msg {
				return listTablesMsg{
					conn:     node.conn,
					database: node.Parent.Label,
					node:     node.Children[0],
				}
			}
		}
	}
	return nil
}

func (m *Model) UpdateLayout(width, height int) {
	if !m.ready {
		m.viewport = viewport.New(width, height)
		m.viewport.KeyMap = viewport.KeyMap{}
		m.viewport.SetContent(m.view())
		m.viewport.Height = height - 1
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.Height = height - 1
	}
}

func (m *Model) updateView(flatten bool) {
	if flatten {
		m.FlattenNodes = flattenAll(m.Roots, m.query)
	}
	m.viewport.SetContent(m.view())
}

func (m Model) view() string {
	var b strings.Builder
	for i, n := range m.FlattenNodes {
		style := lipgloss.NewStyle()
		prefix := "  "
		if len(n.Children) > 0 {
			if n.Expanded {
				prefix = "▾ "
			} else {
				prefix = "▸ "
			}
		}

		if i == m.Cursor {
			style = components.FocusedText
		}

		padding := ""
		for j := 0; j < indent(n); j++ {
			padding += "  "
		}

		if n.loading {
			b.WriteString(style.Render(fmt.Sprintf("%s%s%c %s", padding, prefix, n.Icon, m.spinner.View())))
		} else {
			b.WriteString(style.Render(fmt.Sprintf("%s%s%c %s", padding, prefix, n.Icon, n.Label)))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m Model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}

	return lipgloss.JoinVertical(
		0,
		m.textInput.View(),
		m.viewport.View(),
	)
}
