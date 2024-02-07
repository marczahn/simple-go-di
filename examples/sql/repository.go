package sql

import (
	"database/sql"
	"github.com/marczahn/simple-go-di/pkg/di"
)

var sqlRepo = di.NewInstance[*SQLRepo]()
var sqlConn = di.NewInstance[*sql.DB]()

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

func (s *SQLRepo) Select() (any, error) {
	rsl, err := s.conn.Query("...")
	if err != nil {
		return nil, err
	}

	return rsl, nil
}
