package main

import (
	"redfox-tech/assistir_filmes/cmd/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Post("/filmes", services.NovoFilme)
	app.Post("/sugerir-filmes", services.SugerirFilmes)
	app.Get("/filmes", services.GetFilmes)
	app.Post("/filmes/assistido/:filmeId", services.MarcarFilmeAssistido)
	app.Delete("/filmes/:filmeId", services.RemoverFilme)

	app.Listen(":8080")
}
