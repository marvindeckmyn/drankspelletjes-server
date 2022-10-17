package game

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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

type GameBody struct {
	GameCategory uuid.UUID         `json:"game_category"`
	Name         map[string]string `json:"name"`
	Alias        map[string]string `json:"alias"`
	Description  map[string]string `json:"description"`
	Highlight    bool              `json:"highlight"`
	Img          string            `json:"img"`
	PlayerCount  int32             `json:"player_count"`
	Order        int32             `json:"order"`
}

type GameURL struct {
	ID uuid.UUID `json:"id"`
}

// validateGameBody checks if the body is valid.
func validateGameBody(requestBody io.Reader) (*GameBody, error) {
	v := validator.V{
		"game_category": validator.IsUUIDV4,
		"name":          validator.IsMapStrStr,
		"alias":         validator.IsMapStrStr,
		"description":   validator.IsMapStrStr,
		"highlight":     validator.IsBool,
		"img":           validator.IsString,
		"player_count":  validator.IsInt,
		"order":         validator.IsInt,
	}

	body := GameBody{}

	err := v.ValidateAndMarshalBody(requestBody, &body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &body, nil
}

// validateGameURL checks if the game URL is valid.
func validateGameURL(r *server.Request) (*GameURL, error) {
	v := validator.V{
		"id": validator.IsUUIDV4,
	}

	url := GameURL{}

	err := v.ValidateAndMarshalURL(r, &url)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &url, nil
}

// GetGamesByCategory to retrieve all the games by category.
func GetGamesByCategory(rw server.ResponseWriter, r *server.Request) {
	// Get category
	url, err := validateGameURL(r)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	category := gameModel.GameCategory{
		ID: &url.ID,
	}

	err = gameDao.GetCategory(&category)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	games, err := gameDao.GetGamesByCategory(&category)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	rw.JSON(http.StatusOK, games)
}

// PostGame inserts a game in the database.
func PostGame(rw server.ResponseWriter, r *server.Request) {
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

	// Validate game body
	body, err := validateGameBody(r.R.Body)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	// base64 to image
	gameUuid := types.Ptr(uuid.UUIDv4())

	if body.Img != "" {
		data, err := base64.StdEncoding.DecodeString(body.Img)
		if err != nil {
			log.Error(err.Error())
			rw.JSON(http.StatusInternalServerError, nil)
			return
		}

		filename := fmt.Sprintf("img/game_%s.png", gameUuid.String())

		body.Img = filename

		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			log.Error(err.Error())
			rw.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	// Add game
	game := gameModel.Game{
		ID:           gameUuid,
		GameCategory: &body.GameCategory,
		Name:         &body.Name,
		Alias:        &body.Alias,
		Description:  &body.Description,
		Highlight:    &body.Highlight,
		Img:          &body.Img,
		PlayerCount:  &body.PlayerCount,
		Order:        &body.Order,
		CreatedAt:    types.Ptr(time.Now().UTC()),
	}

	err = gameDao.InsertGame(&game)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	rw.JSON(http.StatusOK, game)
}
