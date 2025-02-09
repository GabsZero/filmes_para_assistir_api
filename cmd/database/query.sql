-- name: GetFilme :one
SELECT * FROM filmes
WHERE id = $1 LIMIT 1;

-- name: CountFilme :one
SELECT count(*) FROM filmes;

-- name: ListFilmes :many
SELECT * FROM filmes
ORDER BY nome
OFFSET $1
LIMIT $2;

-- name: CreateFilme :one
INSERT INTO filmes (
  nome, tipo_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateFilme :exec
UPDATE filmes
  set nome = $2,
  tipo_id = $3,
  assistido = $4
WHERE id = $1;

-- name: DeleteFilme :exec
DELETE FROM filmes
WHERE id = $1;