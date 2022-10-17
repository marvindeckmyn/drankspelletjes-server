package gameDao

import (
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	"github.com/marvindeckmyn/drankspelletjes-server/dao"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	gameModel "github.com/marvindeckmyn/drankspelletjes-server/model/game"
)

var colNamesGame = map[string]string{
	"ID":           "id",
	"GameCategory": "game_category",
	"Name":         "name",
	"Alias":        "alias",
	"PlayerCount":  "player_count",
	"Img":          "img",
	"Description":  "description",
	"Highlight":    "highlight",
	"Views":        "views",
	"Order":        `"order"`,
	"CreatedAt":    "created_at",
}

// unmarshalGame parses the database row to the game object.
func unmarshalGame(game *gameModel.Game, r cdb.CdbResult) error {
	r.UUID("id", &game.ID)
	r.UUID("game_category", &game.GameCategory)
	r.MapStrStr("name", &game.Name)
	r.MapStrStr("alias", &game.Alias)
	r.Int32("player_count", &game.PlayerCount)
	r.Str("img", &game.Img)
	r.MapStrStr("description", &game.Description)
	r.Bool("highlight", &game.Highlight)
	r.Int32("order", &game.Order)

	if r.HasErrorsLog("unmarshal game", "") {
		return &cdb.ErrParseResult{}
	}

	return nil
}

// GetGamesByCategory fetches all the games by category.
func GetGamesByCategory(category *gameModel.GameCategory) ([]*gameModel.Game, error) {
	games := []*gameModel.Game{}

	stmt := cdb.Prepare(`
		select id, game_category, name, alias,
			player_count, img,
			description, highlight, views, "order"
		from game
		where game_category = :category:
		order by "order"
	`)

	stmt.Bind("category", *category.ID)

	rows, err := dao.ExecuteStmt(stmt)
	if err != nil {
		log.Error(err.Error())
		return games, err
	}

	for _, rowGame := range rows {
		game := gameModel.Game{}

		err = unmarshalGame(&game, rowGame)
		if err != nil {
			log.Error(err.Error())
			return []*gameModel.Game{}, err
		}

		games = append(games, &game)
	}

	return games, err
}

// GetGame fetches the game that matches with the non nil values from the given game.
func GetGame(game *gameModel.Game) error {
	fields := cdb.CreateFields(colNamesGame)
	stmt := cdb.PrepareSelect("game", fields, "game", colNamesGame, game)
	rows, err := cdb.Exec(&stmt)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return unmarshalGame(game, rows[0])
}

// InsertGame inserts the game in the database
func InsertGame(game *gameModel.Game) error {
	stmt, err := cdb.PrepareInsert("game", colNamesGame, game)
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

// UpdateGame updates the given game in the database.
func UpdateGame(game *gameModel.Game,
	selectors map[string]interface{}) error {

	stmt, err := cdb.PrepareUpdate("game", colNamesGame, game, selectors)
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

// DeleteGame deletes the given game in the database.
func DeleteGame(game *gameModel.Game) error {
	stmt := cdb.PrepareDelete("game", colNamesGame, game)
	_, err := cdb.Exec(&stmt)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
