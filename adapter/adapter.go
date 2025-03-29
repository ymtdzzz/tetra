package adapter

type Status struct {
	Opened         bool
	DatabaseLoaded bool
	TableLoaded    map[string]bool
}

type ColumnInfo struct {
	Name     string `db:"COLUMN_NAME"`
	Type     string `db:"COLUMN_TYPE"`
	Nullable string `db:"IS_NULLABLE"`
	Key      string `db:"COLUMN_KEY"`
	Extra    string `db:"EXTRA"`
}

type ConstraintInfo struct {
	Name string `db:"CONSTRAINT_NAME"`
	Type string `db:"CONSTRAINT_TYPE"`
}

type ForeignKeyInfo struct {
	ConstraintName       string `db:"CONSTRAINT_NAME"`
	ColumnName           string `db:"COLUMN_NAME"`
	ReferencedTableName  string `db:"REFERENCED_TABLE_NAME"`
	ReferencedColumnName string `db:"REFERENCED_COLUMN_NAME"`
}

type IndexInfo struct {
	IndexName  string `db:"INDEX_NAME"`
	ColumnName string `db:"COLUMN_NAME"`
	NonUnique  int    `db:"NON_UNIQUE"` // 0 = UNIQUE, 1 = NON-UNIQUE

}

type Adapter interface {
	Status() *Status
	Open() error
	Close() error
	ListDatabases() ([]string, error)
	ListTables(database string) ([]string, error)
}
