package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
)

var secretKey = []byte("drankspelletjes_key")

// createToken creates a JWT token for the given account.
func createToken(acc *accountModel.Account) (string, error) {
	if acc == nil || acc.ID == nil {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": fmt.Sprint(acc.ID),
	})

	jwt, err := token.SignedString(secretKey)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	return jwt, err
}
