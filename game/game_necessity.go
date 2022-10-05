package game

import (
	"io"
	"net/http"

	"github.com/marvindeckmyn/drankspelletjes-server/auth"
	accountDao "github.com/marvindeckmyn/drankspelletjes-server/dao/account"
	gameDao "github.com/marvindeckmyn/drankspelletjes-server/dao/game"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
	gameModel "github.com/marvindeckmyn/drankspelletjes-server/model/game"
	"github.com/marvindeckmyn/drankspelletjes-server/server"
	"github.com/marvindeckmyn/drankspelletjes-server/types"
	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
	"github.com/marvindeckmyn/drankspelletjes-server/validator"
)

type GameNecessityBody struct {
	Game uuid.UUID         `json:"game"`
	Name map[string]string `json:"name"`
}

// validateGameNecessityBody checks if the body is valid.
func validateGameNecessityBody(requestBody io.Reader) (*GameNecessityBody, error) {
	v := validator.V{
		"game": validator.IsUUIDV4,
		"name": validator.IsMapStrStr,
	}

	body := GameNecessityBody{}

	err := v.ValidateAndMarshalBody(requestBody, &body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &body, nil
}

// PostGameNecessity inserts a game necessity in the database.
func PostGameNecessity(rw server.ResponseWriter, r *server.Request) {
	// Check account
	accID, err := auth.GetID(r)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusUnauthorized, nil)
		return
	}

	acc := accountModel.Account{
		ID: &accID,
	}

	err = accountDao.GetAccount(&acc)
	if err != nil {
		rw.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Validate game necessity body
	body, err := validateGameNecessityBody(r.R.Body)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	// Add game necessity
	gameNecessity := gameModel.GameNecessity{
		ID:   types.Ptr(uuid.UUIDv4()),
		Game: &body.Game,
		Name: &body.Name,
	}

	err = gameDao.InsertNecessity(&gameNecessity)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	rw.JSON(http.StatusOK, gameNecessity)
}
