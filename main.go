package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marvindeckmyn/drankspelletjes-server/auth"
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
)

// initDB initializes the database.
func initDB() {
	err := cdb.Init("localhost", 5432, "postgres", "aarsaars", "drankspelletjes")
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

// main executes the main function.
func main() {
	r := gin.Default()
	initDB()

	//r.POST("/api/auth/register", auth.Register)
	r.POST("/api/auth/login", auth.Login)

	log.Info("Starting on 1337")
	err := r.Run(":1337")
	if err != nil {
		panic(err)
	}
}
