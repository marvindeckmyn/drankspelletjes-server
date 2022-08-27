package gameDao

import (
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	"github.com/marvindeckmyn/drankspelletjes-server/dao"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	gameModel "github.com/marvindeckmyn/drankspelletjes-server/model/game"
)

var colNamesCategory = map[string]string{
	"ID":    "id",
	"Name":  "name",
	"Order": `"order"`,
}

// unmarshalCategory parses the database row to the category object.
func unmarshalCategory(category *gameModel.GameCategory, r cdb.CdbResult) error {
	r.UUID("id", &category.ID)
	r.MapStrStr("name", &category.Name)
	r.Int32("order", &category.Order)

	if r.HasErrorsLog("unmarshal category", "") {
		return &cdb.ErrParseResult{}
	}

	return nil
}

// GetCategories fetches all the categories.
func GetCategories() ([]*gameModel.GameCategory, error) {
	categories := []*gameModel.GameCategory{}

	stmt := cdb.Prepare(`
		select id, name, "order"
		from game_category
		order by "order"
	`)

	rows, err := dao.ExecuteStmt(stmt)
	if err != nil {
		log.Error(err.Error())
		return categories, err
	}

	for _, rowCategory := range rows {
		category := gameModel.GameCategory{}

		err = unmarshalCategory(&category, rowCategory)
		if err != nil {
			log.Error(err.Error())
			return []*gameModel.GameCategory{}, err
		}

		categories = append(categories, &category)
	}

	return categories, err
}

// GetCategory fetches the category that matches with the non nil values from the given category.
func GetCategory(category *gameModel.GameCategory) error {
	fields := cdb.CreateFields(colNamesCategory)
	stmt := cdb.PrepareSelect("game_category", fields, "gc", colNamesCategory, category)
	rows, err := cdb.Exec(&stmt)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return unmarshalCategory(category, rows[0])
}

// InsertCategory inserts the category in the database
func InsertCategory(category *gameModel.GameCategory) error {
	stmt, err := cdb.PrepareInsert("game_category", colNamesCategory, category)
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

// UpdateCategory updates the given category in the database.
func UpdateCategory(category *gameModel.GameCategory,
	selectors map[string]interface{}) error {

	stmt, err := cdb.PrepareUpdate("game_category", colNamesCategory, category, selectors)
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

// DeleteCategory deletes the given category in the database.
func DeleteCategory(category *gameModel.GameCategory) error {
	stmt := cdb.PrepareDelete("game_category", colNamesCategory, category)
	_, err := cdb.Exec(&stmt)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
