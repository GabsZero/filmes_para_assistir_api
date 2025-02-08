package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func main() {

	ctx := context.Background()

	// só pra fazer funcionar durante desenvolvimento
	conn, err := pgx.Connect(ctx, "user=root password=postgres dbname=assistir_filmes_api host=filmes_para_assistir_db port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":8080")
}
