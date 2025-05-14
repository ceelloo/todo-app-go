package main

import (
	"log"

	"github.com/ceelloo/todo-app-go/internal/database"
)


func main() {
	db := database.New()
	database.InitializeDatabase(db)

	app := application{
		addr: ":3000",
		env: "prod",
		db: db,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
