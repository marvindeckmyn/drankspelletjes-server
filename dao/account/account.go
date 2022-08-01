package accountModel

import (
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	"github.com/marvindeckmyn/drankspelletjes-server/dao"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
)

// a mapping from the model names to the db names
var colNamesAccount = map[string]string{
	"ID":       "id",
	"Name":     "name",
	"Email":    "email",
	"Password": "password",
}

// unmarshalAccount parses the database row to the account object.
func unmarshalAccount(acc *accountModel.Account, r cdb.CdbResult) error {
	r.UUID("id", &acc.ID)
	r.Str("name", &acc.Name)
	r.Str("email", &acc.Email)
	r.Str("password", &acc.Password)

	if r.HasErrorsLog("unmarshal account", "") {
		return &cdb.ErrParseResult{}
	}

	return nil
}

// GetAccount fetches the account that matches with the non nil values from the given account.
func GetAccount(acc *accountModel.Account) error {
	fields := cdb.CreateFields(colNamesAccount)
	stmt := cdb.PrepareSelect("account", fields, "a", colNamesAccount, acc)
	rows, err := dao.ExecuteStmt(stmt)
	if err != nil {
		log.Error("Account not found")
		return &cdb.ErrQuery{Cause: err}
	}

	if len(rows) == 0 {
		return &cdb.ErrMissingResult{}
	}

	return unmarshalAccount(acc, rows[0])
}

// InsertAccount inserts the account in the database.
func InsertAccount(acc *accountModel.Account) error {
	stmt, err := cdb.PrepareInsert("account", colNamesAccount, acc)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	_, err = cdb.Exec(&stmt)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
