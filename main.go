package main

import (
	"github.com/marvindeckmyn/drankspelletjes-server/account"
	"github.com/marvindeckmyn/drankspelletjes-server/auth"
	"github.com/marvindeckmyn/drankspelletjes-server/cdb"
	"github.com/marvindeckmyn/drankspelletjes-server/game"
	"github.com/marvindeckmyn/drankspelletjes-server/log"
	"github.com/marvindeckmyn/drankspelletjes-server/server"
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
	s := server.New()
	initDB()

	s.Get("/api/auth/account", account.Get)

	//s.Post("/api/auth/register", auth.Register)
	s.Post("/api/auth/login", auth.Login)
	s.Post("/api/auth/logout", auth.Logout)

	s.Get("/api/category", game.GetCategories)
	s.Get("/api/category/{id}", game.GetCategoryById)
	s.Post("/api/category", game.PostCategory)
	s.Put("/api/category/{id}", game.UpdateCategory)
	s.Delete("/api/category/{id}", game.DeleteCategory)

	log.Info("Starting on 1337")
	err := s.ListenAndServe(1337)
	if err != nil {
		panic(err)
	}
}
