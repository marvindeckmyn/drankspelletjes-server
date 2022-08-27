package game

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marvindeckmyn/drankspelletjes-server/auth"
	accountDao "github.com/marvindeckmyn/drankspelletjes-server/dao/account"
	gameDao "github.com/marvindeckmyn/drankspelletjes-server/dao/game"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
	gameModel "github.com/marvindeckmyn/drankspelletjes-server/model/game"
	"github.com/marvindeckmyn/drankspelletjes-server/types"
	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
	"github.com/marvindeckmyn/drankspelletjes-server/validator"
)

type CategoryBody struct {
	Name  map[string]string `json:"name"`
	Order int32             `json:"order"`
}

type CategoryURL struct {
	ID uuid.UUID `json:"id"`
}

// validateCategoryBody checks if the body is valid.
func validateCategoryBody(requestBody io.Reader) (*CategoryBody, error) {
	v := validator.V{
		"name":  validator.IsMapStrStr,
		"order": validator.IsInt,
	}

	body := CategoryBody{}

	err := v.ValidateAndMarshalBody(requestBody, &body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &body, nil
}

// validateCategoryURL checks if the product categogry URL is valid.
func validateCategoryURL(c *gin.Context) (*CategoryURL, error) {
	v := validator.V{
		"id": validator.IsUUIDV4,
	}

	url := CategoryURL{}

	err := v.ValidateAndMarshalURL(c.Request, &url)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &url, nil
}

// GetCategories to retrieve all the categories.
func GetCategories(c *gin.Context) {
	categories, err := gameDao.GetCategories()
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryById to retrieve a category by UUID.
func GetCategoryById(c *gin.Context) {
	// Get category
	url, err := validateCategoryURL(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	category := gameModel.GameCategory{
		ID: &url.ID,
	}

	err = gameDao.GetCategory(&category)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, category)
}

// PostCategory inserts a category in the database.
func PostCategory(c *gin.Context) {
	// Check account
	accID, err := auth.GetID(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	acc := accountModel.Account{
		ID: &accID,
	}

	err = accountDao.GetAccount(&acc)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Validate category body
	body, err := validateCategoryBody(c.Request.Body)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Add category
	category := gameModel.GameCategory{
		ID:    types.Ptr(uuid.UUIDv4()),
		Name:  &body.Name,
		Order: &body.Order,
	}

	err = gameDao.InsertCategory(&category)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory updates a selected category the database.
func UpdateCategory(c *gin.Context) {
	// Check account
	accID, err := auth.GetID(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	acc := accountModel.Account{
		ID: &accID,
	}

	err = accountDao.GetAccount(&acc)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Validate category URL
	url, err := validateCategoryURL(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Validate category body
	body, err := validateCategoryBody(c.Request.Body)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Get category
	category := gameModel.GameCategory{
		ID: &url.ID,
	}

	err = gameDao.GetCategory(&category)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Update category
	category.Name = &body.Name
	category.Order = &body.Order

	selectors := map[string]interface{}{
		"ID": *&category.ID,
	}

	err = gameDao.UpdateCategory(&category, selectors)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeletCategory deletes a category in the database.
func DeleteCategory(c *gin.Context) {
	// Check account
	accID, err := auth.GetID(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	acc := accountModel.Account{
		ID: &accID,
	}

	err = accountDao.GetAccount(&acc)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Validate category URL
	url, err := validateCategoryURL(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Get category
	category := gameModel.GameCategory{
		ID: &url.ID,
	}

	err = gameDao.GetCategory(&category)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Delete category
	err = gameDao.DeleteCategory(&category)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}
