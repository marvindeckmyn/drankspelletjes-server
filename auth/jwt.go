package auth

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
	"github.com/marvindeckmyn/drankspelletjes-server/server"
	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
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

// ParseToken to parse a JWT token and make it readable
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, nil
}

// GetID to get the ID of a JWT
func GetID(r *server.Request) (uuid.UUID, error) {
	jwt, err := r.Cookie("drnkngg-token")
	if err != nil {
		return uuid.UUID{}, err
	}

	parsedToken, err := ParseToken(*jwt)
	if err != nil {
		return uuid.UUID{}, err
	}

	id := parsedToken["ID"]

	marshalToken, err := json.Marshal(id)
	if err != nil {
		return uuid.UUID{}, err
	}

	token := uuid.UUID{}
	err = json.Unmarshal(marshalToken, &token)
	if err != nil {
		return uuid.UUID{}, err
	}

	return token, nil
}
