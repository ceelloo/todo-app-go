package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	Id        int    `json:"_id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func (app *application) createTodo(c *fiber.Ctx) error {
	todo := Todo{}

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Body is required"})
	}

	result, err := app.db.Exec(`INSERT INTO task (completed, body) VALUES (?, ?)`, todo.Completed, todo.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to insert task"})
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to get last insert id"})
	}

	todo.Id = int(id)
	return c.Status(200).JSON(todo)
}

func (app *application) getTodoFromParams(c *fiber.Ctx) error {
	q, err := app.db.Query("SELECT id, completed, body FROM task WHERE id = ?", c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to get task"})
	}

	defer q.Close()

	t := Todo{}

	if !q.Next() {
		return c.Status(404).JSON(fiber.Map{"message": "Task not found"})
	}

	q.Scan(&t.Id, &t.Completed, &t.Body)
	return c.Status(200).JSON(t)
}

func (app *application) getTodos(c *fiber.Ctx) error {
	todos := []Todo{}

	q, err := app.db.Query("SELECT id, completed, body FROM task")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to get tasks"})
	}

	defer q.Close()

	for q.Next() {
		t := Todo{}
		err := q.Scan(&t.Id, &t.Completed, &t.Body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Failed to scan task"})
		}
		todos = append(todos, t)
	}

	return c.Status(200).JSON(todos)
}

func (app *application) updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	var data struct {
		Completed bool `json:"completed"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	res, err := app.db.Exec("UPDATE task SET completed = ? WHERE id = ?", data.Completed, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to update completed status"})
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to check result"})
	}

	if rows == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "Task not found"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":   "Completed status updated",
		"completed": data.Completed,
	})
}

func (app *application) deleteTodo(c *fiber.Ctx) error {
	id := fmt.Sprint(c.Params("id"))

	res, err := app.db.Exec("DELETE FROM task WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to delete task"})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to check deleted row"})
	}

	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "Task not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "successfully deleted"})
}
