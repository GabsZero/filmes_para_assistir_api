create table if not exists filmes(
    id BIGSERIAL not null primary key,
    nome VARCHAR(255) not null,
    assistido boolean not null default false,
    tipo_id BIGSERIAL not null,
    CONSTRAINT fk_tipo FOREIGN KEY(tipo_id) REFERENCES tipos(id),
    created_at date default CURRENT_TIMESTAMP
);