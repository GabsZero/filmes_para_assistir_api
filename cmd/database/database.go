package database

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5"
)

func InitDB() *Queries {
	ctx := context.Background()

	// só pra fazer funcionar durante desenvolvimento
	conn, err := pgx.Connect(ctx, "user=root password=postgres dbname=assistir_filmes_api host=filmes_para_assistir_db port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(
		"file://cmd/database/migrations",
		"postgres://filmes_para_assistir_db:5432/assistir_filmes_api?sslmode=enable")
	m.Steps(2)

	if err != nil {
		log.Fatal(err)
	}

	// defer conn.Close(ctx)

	return New(conn)
}
