// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"
)

const createFilme = `-- name: CreateFilme :one
INSERT INTO filmes (
  nome, tipo_id
) VALUES (
  $1, $2
)
RETURNING id, nome, assistido, tipo_id, created_at
`

type CreateFilmeParams struct {
	Nome   string
	TipoID int64
}

func (q *Queries) CreateFilme(ctx context.Context, arg CreateFilmeParams) (Filme, error) {
	row := q.db.QueryRow(ctx, createFilme, arg.Nome, arg.TipoID)
	var i Filme
	err := row.Scan(
		&i.ID,
		&i.Nome,
		&i.Assistido,
		&i.TipoID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteFilme = `-- name: DeleteFilme :exec
DELETE FROM filmes
WHERE id = $1
`

func (q *Queries) DeleteFilme(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteFilme, id)
	return err
}

const getFilme = `-- name: GetFilme :one
SELECT id, nome, assistido, tipo_id, created_at FROM filmes
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFilme(ctx context.Context, id int64) (Filme, error) {
	row := q.db.QueryRow(ctx, getFilme, id)
	var i Filme
	err := row.Scan(
		&i.ID,
		&i.Nome,
		&i.Assistido,
		&i.TipoID,
		&i.CreatedAt,
	)
	return i, err
}

const listFilmes = `-- name: ListFilmes :many
SELECT id, nome, assistido, tipo_id, created_at FROM filmes
ORDER BY nome
`

func (q *Queries) ListFilmes(ctx context.Context) ([]Filme, error) {
	rows, err := q.db.Query(ctx, listFilmes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Filme
	for rows.Next() {
		var i Filme
		if err := rows.Scan(
			&i.ID,
			&i.Nome,
			&i.Assistido,
			&i.TipoID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFilme = `-- name: UpdateFilme :exec
UPDATE filmes
  set nome = $2,
  assistido = $3
WHERE id = $1
`

type UpdateFilmeParams struct {
	ID        int64
	Nome      string
	Assistido bool
}

func (q *Queries) UpdateFilme(ctx context.Context, arg UpdateFilmeParams) error {
	_, err := q.db.Exec(ctx, updateFilme, arg.ID, arg.Nome, arg.Assistido)
	return err
}
