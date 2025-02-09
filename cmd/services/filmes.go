package services

import (
	"context"
	"fmt"
	"log"
	"redfox-tech/assistir_filmes/cmd/database"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type JsonResponse struct {
	Data   interface{} `json:"data"`
	Error  string      `json:"error,omitempty"`
	Status int         `json:"status"`
}

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type NovoFilmeDto struct {
	Nome   string `validate:"required" json:"nome" form:"nome"`
	TipoID string `json:"tipoID" form:"tipoID"`
}

func MarcarFilmeAssistido(c *fiber.Ctx) error {
	filmeId := c.Params("filmeId")
	filmeIdInt, err := strconv.Atoi(filmeId)

	if err != nil {
		c.Status(400)
		return c.JSON(JsonResponse{
			Data:   nil,
			Error:  "ID do filme inválido",
			Status: 400,
		})
	}

	db := database.InitDB()

	filme, err := db.GetFilme(context.Background(), int64(filmeIdInt))
	if err != nil {
		c.Status(404)
		return c.JSON(JsonResponse{
			Data:   nil,
			Error:  "Filme não encontrado",
			Status: 404,
		})
	}

	filme.Assistido = true

	db.UpdateFilme(context.Background(), database.UpdateFilmeParams{
		ID:        filme.ID,
		Nome:      filme.Nome,
		TipoID:    filme.TipoID,
		Assistido: filme.Assistido,
	})

	return c.JSON(JsonResponse{
		Data:   filme,
		Error:  "",
		Status: 200,
	})
}

func NovoFilme(c *fiber.Ctx) error {
	novoFilme := NovoFilmeDto{}
	fmt.Println(c.Request())
	c.BodyParser(&novoFilme)

	validate := validator.New()

	errors := validate.Struct(novoFilme)

	if errors != nil {
		c.Status(400)
		validationErrors := []string{}
		for _, err := range errors.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf(
				"O campo [%s]: '%v' | é '%s'",
				err.Field(),
				err.Value(),
				err.Tag(),
			))
		}
		return c.JSON(JsonResponse{
			Data:   nil,
			Error:  strings.Join(validationErrors, ", "),
			Status: 400,
		})
	}

	tipoId, err := strconv.Atoi(novoFilme.TipoID)

	if err != nil {
		log.Fatal(err)
	}

	db := database.InitDB()

	filme, err := db.CreateFilme(c.Context(), database.CreateFilmeParams{
		Nome:   novoFilme.Nome,
		TipoID: int64(tipoId),
	})

	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(JsonResponse{
		Data:   filme,
		Error:  "",
		Status: 200,
	})
}
