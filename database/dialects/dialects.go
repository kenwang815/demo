package dialects

type Dialect string

func (d Dialect) String() string {
	return string(d)
}

const (
	Postgres Dialect = "postgres"
	Sqlite   Dialect = "sqlite"
)
