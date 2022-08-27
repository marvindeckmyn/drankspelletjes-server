package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marvindeckmyn/drankspelletjes-server/account"
	"github.com/marvindeckmyn/drankspelletjes-server/auth"
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	"github.com/marvindeckmyn/drankspelletjes-server/game"
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

	r.GET("/api/auth/account", account.Get)

	//r.POST("/api/auth/register", auth.Register)
	r.POST("/api/auth/login", auth.Login)
	r.POST("/api/auth/logout", auth.Logout)

	r.GET("/api/category", game.GetCategories)
	r.GET("/api/category/:id", game.GetCategoryById)
	r.POST("/api/category", game.PostCategory)
	r.PUT("/api/category/:id", game.UpdateCategory)
	r.DELETE("/api/category/:id", game.DeleteCategory)

	log.Info("Starting on 1337")
	err := r.Run(":1337")
	if err != nil {
		panic(err)
	}
}
