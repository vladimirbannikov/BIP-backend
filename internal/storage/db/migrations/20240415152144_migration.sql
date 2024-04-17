-- +goose Up
-- +goose StatementBegin
create schema if not exists users_schema;

create table if not exists users_schema.users (
    id SERIAL PRIMARY KEY NOT NULL,
    login text unique not null,
    password_hash text not null,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

create schema if not exists auth_schema;

create table if not exists auth_schema.users_secrets (
    user_id SERIAL not null references users_schema.users(id),
    secret text,
    session_id text,
    unique (user_id, session_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema users_schema cascade;
drop schema auth_schema cascade;
-- +goose StatementEnd
