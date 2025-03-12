-- name: GetFilme :one
SELECT * FROM filmes
WHERE id = $1 LIMIT 1;

-- name: CountFilme :one
SELECT count(*) FROM filmes;

-- name: ListFilmes :many
SELECT f.id, f.nome, f.assistido, t.nome as "tipo" FROM filmes f
join tipos t on f.tipo_id = t.id
WHERE assistido = $1 
ORDER BY f.nome
OFFSET $2
LIMIT $3;

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