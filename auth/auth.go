package auth

import (
	"net/http"
	"strings"

	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	accountDao "github.com/marvindeckmyn/drankspelletjes-server/dao/account"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	accountModel "github.com/marvindeckmyn/drankspelletjes-server/model/account"
	"github.com/marvindeckmyn/drankspelletjes-server/server"
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
func Register(rw server.ResponseWriter, r *server.Request) {
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

	err := v.ValidateAndMarshalBody(r.R.Body, &body)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
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
			rw.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	log.Info("%s is registering", body.Email)

	// Insert account
	hash, err := hashPassword(body.Password)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusInternalServerError, nil)
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
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	rw.JSON(http.StatusCreated, nil)
}

// checkPasswordHash will check if the password is equal to the hash if hashed.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login to log in an account.
func Login(rw server.ResponseWriter, r *server.Request) {
	// Check body
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	v := validator.V{
		"email":    validator.IsEmail,
		"password": validator.IsString,
	}

	err := v.ValidateAndMarshalBody(r.R.Body, &body)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
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
			rw.JSON(http.StatusNotFound, nil)
			return
		}

		rw.JSON(http.StatusInternalServerError, nil)
		return
	}

	match := checkPasswordHash(body.Password, *acc.Password)
	if !match {
		rw.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Add JWT
	jwt, err := createToken(&acc)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusInternalServerError, nil)
		return
	}

	cookie := &http.Cookie{
		Name:   "drnkngg-token",
		Value:  jwt,
		Path:   "/",
		Domain: ".drankspelletjes.local",
		MaxAge: 3600 * 24 * 7,
	}

	http.SetCookie(rw.W, cookie)

	rw.JSON(http.StatusOK, nil)
}

// Logout to log out of an account.
func Logout(rw server.ResponseWriter, r *server.Request) {
	accID, err := GetID(r)
	if err != nil {
		log.Error(err.Error())
		rw.JSON(http.StatusBadRequest, nil)
		return
	}

	log.Info("%s is logging out", accID)

	// Get account
	acc := accountModel.Account{
		ID: &accID,
	}

	err = accountDao.GetAccount(&acc)
	if err != nil {
		log.Error(err.Error())
		if _, ok := err.(*cdb.ErrMissingResult); !ok {
			rw.JSON(http.StatusNotFound, nil)
			return
		}

		rw.JSON(http.StatusInternalServerError, nil)
		return
	}

	// Remove cookie
	cookie := &http.Cookie{
		Name:   "drnkngg-token",
		Value:  "",
		Path:   "/",
		Domain: ".drankspelletjes.local",
		MaxAge: -1,
	}

	http.SetCookie(rw.W, cookie)

	rw.JSON(http.StatusOK, nil)
}
