package dao

import "github.com/marvindeckmyn/drankspelletjes-server/cdb"

// ExecuteStmt executes the given statement and gives rows back
func ExecuteStmt(stmt cdb.Statement) ([]cdb.CdbResult, error) {
	rows, err := cdb.Exec(&stmt)
	if err != nil {
		err = &cdb.ErrQuery{Cause: err}
	}

	if len(rows) == 0 {
		err = &cdb.ErrMissingResult{}
	}

	return rows, err
}
