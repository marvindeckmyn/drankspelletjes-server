package main

import (
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
	initDB()
}
