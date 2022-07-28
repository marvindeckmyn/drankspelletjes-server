package cdb

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"

	"github.com/jackc/pgx/v4"
)

var tx *Transaction = nil

type Transaction struct {
	stmts []Statement
}

func (t *Transaction) AddStmt(s Statement) {
	t.stmts = append(t.stmts, s)
}

func (t *Transaction) Exec() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	con, err := getCon(ctx)
	if err != nil {
		return err
	}

	defer con.Release()

	return crdbpgx.ExecuteTx(ctx, con, pgx.TxOptions{}, func(tx pgx.Tx) error {
		for _, stmt := range t.stmts {
			str, values, err := stmt.getPgQuery()
			if err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, str, values...); err != nil {
				fmt.Println(str, values)
				return err
			}
		}

		return nil
	})
}

func BeginTx() *Transaction {
	tx = &Transaction{
		stmts: []Statement{},
	}

	return tx
}
