package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marvindeckmyn/drankspelletjes-server/auth"
	accountDao "github.com/marvindeckmyn/drankspelletjes-server/dao/account"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
)

// Get to retrieve the account of the current user.
func Get(c *gin.Context) {
	accID, err := auth.GetID(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	account := accountModel.Account{
		ID: &accID,
	}

	err = accountDao.GetAccount(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"name": *account.Name,
	})
}
