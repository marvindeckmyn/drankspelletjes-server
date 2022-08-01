package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	accountDao "github.com/marvindeckmyn/drankspelletjes-server/dao/account"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
	"github.com/marvindeckmyn/drankspelletjes-server/types"
	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
	"github.com/marvindeckmyn/drankspelletjes-server/validator"
	"golang.org/x/crypto/bcrypt"
)

const SALT_ROUNDS = 14

// hashPassword hashes the given password with bcrypt.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), SALT_ROUNDS)
	return string(bytes), err
}

// Register to create an account.
func Register(c *gin.Context) {
	// Check body
	body := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	// Validate type
	v := validator.V{
		"name":     validator.IsString,
		"email":    validator.IsEmail,
		"password": validator.IsString,
	}

	err := v.ValidateAndMarshalBody(c.Request.Body, &body)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Validate email unique
	acc := accountModel.Account{
		Email: &body.Email,
	}

	err = accountDao.GetAccount(&acc)
	if err != nil {
		if !strings.Contains(err.Error(), "No results") {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	log.Info("%s is registering", body.Email)

	// Insert account
	hash, err := hashPassword(body.Password)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	acc = accountModel.Account{
		ID:       types.Ptr(uuid.UUIDv4()),
		Name:     &body.Name,
		Email:    &body.Email,
		Password: &hash,
	}

	err = accountDao.InsertAccount(&acc)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// checkPasswordHash will check if the password is equal to the hash if hashed.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login to log in an account.
func Login(c *gin.Context) {
	// Check body
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	v := validator.V{
		"email":    validator.IsEmail,
		"password": validator.IsString,
	}

	err := v.ValidateAndMarshalBody(c.Request.Body, &body)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	log.Info("%s is logging in", body.Email)

	// Get account
	acc := accountModel.Account{
		Email: &body.Email,
	}

	err = accountDao.GetAccount(&acc)
	if err != nil {
		log.Error("Account not found with email %s", body.Email)
		if _, ok := err.(*cdb.ErrMissingResult); !ok {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	match := checkPasswordHash(body.Password, *acc.Password)
	if !match {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Add JWT
	jwt, err := createToken(&acc)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	// Add cookie with JWT in it
	c.SetCookie("drnkngg-token", jwt, 3600*24*7, "/", ".drankspelletjes.local", false, false)

	c.JSON(http.StatusOK, nil)
}
