package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type application struct {
	addr string
	env  string
	db   *sql.DB
}

func (app *application) mount() *fiber.App {
	init := fiber.New()
	r := init.Group("/api")

	r.Use(logger.New())

	post := r.Group("/task")

	post.Get("/", app.getTodos)
	post.Get("/:id", app.getTodoFromParams)
	post.Post("/", app.createTodo)
	post.Patch("/:id", app.updateTodo)
	post.Delete("/:id", app.deleteTodo)

	return init
}

func (app *application) run(fb *fiber.App) error {
	if app.env == "prod" {
		fb.Static("/", "./web/dist")
	}

	return fb.Listen(app.addr)
}
