package tree

import "github.com/ymtdzzz/tetra/adapter"

type connectionOpenMsg struct {
	conn *adapter.DBConnection
	node *TreeNode
}

type connectionOpenDoneMsg struct {
	conn *adapter.DBConnection
	node *TreeNode
}

type listDatabasesMsg struct {
	conn *adapter.DBConnection
	node *TreeNode
}

type listDatabasesDoneMsg struct {
	conn      *adapter.DBConnection
	databases []string
	node      *TreeNode
}

type listTablesMsg struct {
	conn     *adapter.DBConnection
	database string
	node     *TreeNode
}

type listTablesDoneMsg struct {
	conn   *adapter.DBConnection
	tables []string
	node   *TreeNode
}
