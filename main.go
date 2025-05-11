package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	Id        int    `json:"_id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	initDatabase()
	defer db.Close()

	app := fiber.New()
	api := app.Group("/api/task")

	app.Static("/", "./client/dist")

	api.Get("/", getTodos)
	api.Get("/:id", getTodoFromParams)
	api.Post("/", createTodo)
	api.Patch("/:id", updateTodo)
	api.Delete("/:id", deleteTodo)

	log.Fatal(app.Listen(":3000"))
}
