package sql

import (
	"database/sql"
	"github.com/marczahn/simple-go-di/pkg/di"
)

var sqlRepo = di.NewSingleton[*SQLRepo]()
var sqlConn = di.NewSingleton[*sql.DB]()

func SQLConn() *sql.DB {
	return sqlConn.GetOrSet(
		func() *sql.DB {
			db, err := sql.Open("...", "...")
			if err != nil {
				panic("could not open connection")
			}

			return db
		},
		false,
	)
}

func Repo() *SQLRepo {
	return sqlRepo.GetOrSet(
		func() *SQLRepo {
			return &SQLRepo{conn: SQLConn()}
		},
		false,
	)
}

type SQLRepo struct {
	conn *sql.DB
}

func (s *SQLRepo) Select() {
	s.conn.Query("...")

	// return result
}
