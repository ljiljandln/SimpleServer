package main

import (
	"l0/internal/creatingSchemas"
	"l0/internal/database"
)

func main() {
	db := database.SetConfig()

	db.Open()
	defer db.Close()
	if err := creatingSchemas.CreateDbSchemas(db); err != nil {
		panic(err)
	}
}
