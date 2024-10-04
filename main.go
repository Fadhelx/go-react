package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loaing .env file")
	}
	PORT := os.Getenv("PORT")

	var todos []Todo

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hello World"})
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} // {id: 0, completed: false, body: ""}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body can't be empty, It's required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Todo with id: %s not found!", id)})
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		if len(todos) > 0 {
			return c.Status(200).JSON(fiber.Map{"todos": todos})
		}

		return c.Status(404).JSON(fiber.Map{"error": "No todos found!"})
	})

	app.Delete("api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"message": "Todo deleted successfully"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Todo with id: %s not found!", id)})
	})

	log.Fatal(app.Listen(":" + PORT))
}
