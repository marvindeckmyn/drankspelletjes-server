package gameDao

import (
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	gameModel "github.com/marvindeckmyn/drankspelletjes-server/model/game"
)

var colNamesGameNecessity = map[string]string{
	"ID":   "id",
	"Game": "game",
	"Name": "name",
}

// InsertNecessity inserts the game necessity in the database.
func InsertNecessity(necessity *gameModel.GameNecessity) error {
	stmt, err := cdb.PrepareInsert("game_necessity", colNamesGameNecessity, necessity)
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
