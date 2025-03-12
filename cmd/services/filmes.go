package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"redfox-tech/assistir_filmes/cmd/database"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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

type SugestaoFilmeDto struct {
	Genero string `validate:"required" json:"genero" form:"genero"`
}

type FilmeSugerido struct {
	Nome   string
	Genero string
	Tipo   string
}

func SugerirFilmes(c *fiber.Ctx) error {
	ctx := context.Background()
	sugestaoFilmeDto := SugestaoFilmeDto{}
	c.BodyParser(&sugestaoFilmeDto)

	fmt.Println(sugestaoFilmeDto)

	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GOOGLE_AI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"

	prompt := fmt.Sprintf(`Sugira alguns filmes e séries populares com base no gênero %s e informe se é filme ou série no campo tipo usando o seguinte JSON schema:
                   Filme = {'nome': string, 'genero': string, 'tipo': string}
	           Return: Array<Filme>`, sugestaoFilmeDto.Genero)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		sb.WriteString(fmt.Sprintf("%v ", part)) // Se precisar de algo específico, adapte aqui
	}

	filmesSugeridos := []FilmeSugerido{}
	fmt.Println(json.Unmarshal([]byte(sb.String()), &filmesSugeridos))

	fmt.Println(filmesSugeridos)

	return c.JSON(JsonResponse{
		Data:   filmesSugeridos,
		Error:  "",
		Status: 200,
	})
}

func GetFilmes(c *fiber.Ctx) error {
	db := database.InitDB()

	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("page", "1")
	filtroAssistido := c.Query("assistido", "false")

	assistido, _ := strconv.ParseBool(filtroAssistido)

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	page := offset

	if offset == 1 {
		offset = 0
	} else {
		offset = (offset - 1) * limit
	}

	filmes, err := db.ListFilmes(context.Background(), database.ListFilmesParams{
		Assistido: assistido,
		Offset:    int32(offset),
		Limit:     int32(limit),
	})

	if err != nil {
		c.Status(500)
		return c.JSON(JsonResponse{
			Data:   nil,
			Error:  "Não foi possível recuperar os filmes",
			Status: 500,
		})
	}

	total, _ := db.CountFilme(context.Background())
	totalPages := 1
	countingPages := float64(total) / float64(limit)

	mathResult, _ := math.Modf(countingPages)

	if mathResult > 0 {
		totalPages = int(mathResult) + 1
	}

	result := make(map[string]interface{})
	result["filmes"] = filmes
	result["total"] = total
	result["page"] = page
	result["count_current_page"] = len(filmes)
	result["total_pages"] = totalPages

	return c.JSON(JsonResponse{
		Data:   result,
		Error:  "",
		Status: 200,
	})
}

func RemoverFilme(c *fiber.Ctx) error {
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

	erroAoGravar := db.DeleteFilme(context.Background(), int64(filmeIdInt))
	if erroAoGravar != nil {
		c.Status(404)
		return c.JSON(JsonResponse{
			Data:   nil,
			Error:  "Filme não foi removido",
			Status: 404,
		})
	}

	return c.JSON(JsonResponse{
		Data:   nil,
		Error:  "",
		Status: 200,
	})
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
