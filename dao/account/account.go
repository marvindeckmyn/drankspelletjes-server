package accountModel

import (
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
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
