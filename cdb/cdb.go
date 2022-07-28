package cdb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// The maximum duration before cockroach should answer a request.
const requestTimeout time.Duration = 20 * time.Second

// dbPool represents the Connection pool which holds the connections to perform operations to the
// database.
var dbPool *pgxpool.Pool = nil

// Init instantiates a db connection pool.
func Init(host string, port uint16, user string, password string, database string) error {
	url := "postgres://" + user + ":" + password + "@" + host + ":" + fmt.Sprint(port) + "/" + database
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return &ErrConnect{Cause: err}
	}

	dbPool = pool
	return nil
}

// GetCon returns a connection which is obtained from the existing connection pool.
func getCon(ctx context.Context) (*pgxpool.Conn, error) {
	if dbPool == nil {
		return nil, &ErrNotInstantiated{}
	}

	return dbPool.Acquire(ctx)
}

// parseRows transforms the database result into cockroach db results.
func parseRows(rows pgx.Rows) []CdbResult {
	res := []CdbResult{}

	for rows.Next() {
		row := CdbResult{
			data: map[string]interface{}{},
			errs: []error{},
		}

		values, err := rows.Values()
		if err != nil {
			return nil
		}

		for idx, value := range values {
			field := rows.FieldDescriptions()[idx]
			row.data[string(field.Name)] = value
		}

		res = append(res, row)
	}

	return res
}

// Handle validates whether there's an ongoing transaction, if there is one it will add the statement
// to the transaction statements. Otherwise it will call the Exec function.
func Handle(s *Statement) ([]CdbResult, error) {
	if tx != nil {
		tx.AddStmt(*s)
		return nil, nil
	}

	return Exec(s)
}

// Exec executes the statement on the database and returns the results from it.
func Exec(s *Statement) ([]CdbResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	con, err := getCon(ctx)
	if err != nil {
		return nil, &ErrConnect{Cause: err}
	}

	str, values, err := s.getPgQuery()
	if err != nil {
		return nil, &ErrFailedToParseQuery{Cause: err}
	}

	rows, err := con.Query(ctx, str, values...)
	if err != nil {
		return nil, &ErrQuery{err}
	}

	defer con.Release()
	defer rows.Close()

	data := parseRows(rows)

	if rows.Err() != nil {
		return nil, &ErrQuery{Cause: rows.Err()}
	}

	return data, nil
}

// // Exec executes the statement on the database and returns the results from it.
// func Exec(s *Statement) ([]CdbResult, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
// 	defer cancel()

// 	con, err := getCon(ctx)
// 	if err != nil {
// 		return nil, &ErrConnect{Cause: err}
// 	}

// 	str, values, err := s.getPgQuery()
// 	if err != nil {
// 		return nil, &ErrFailedToParseQuery{Cause: err}
// 	}

// 	cmdTag, err := con.Exec(ctx, str, values...)
// 	if err != nil {
// 		return nil, &ErrQuery{err}
// 	}

// 	defer con.Release()
// 	// defer rows.Close()

// 	fmt.Println(cmdTag)

// 	// data := parseRows(rows)
// 	tmp := []CdbResult{}
// 	return tmp, nil
// }
