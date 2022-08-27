package account

import (
	"net/http"

	"github.com/marvindeckmyn/drankspelletjes-server/auth"
	accountDao "github.com/marvindeckmyn/drankspelletjes-server/dao/account"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
	"github.com/marvindeckmyn/drankspelletjes-server/server"
)

// Get to retrieve the account of the current user.
func Get(rw server.ResponseWriter, r *server.Request) {
	accID, err := auth.GetID(r)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusUnauthorized, nil)
		return
	}

	account := accountModel.Account{
		ID: &accID,
	}

	err = accountDao.GetAccount(&account)
	if err != nil {
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	rw.JSON(http.StatusOK, map[string]interface{}{
		"name": *account.Name,
	})
}
