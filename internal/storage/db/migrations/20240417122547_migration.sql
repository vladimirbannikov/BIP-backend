-- +goose Up
-- +goose StatementBegin
create schema if not exists users_schema;

create table if not exists users_schema.users (
    login text PRIMARY KEY not null,
    password_hash text not null,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

create table if not exists users_schema.user_profile (
    login text primary key not null references users_schema.users(login),
    info text
);

create schema if not exists auth_schema;

create table if not exists auth_schema.users_secrets (
    login text not null references users_schema.users(login),
    secret text,
    session_id text,
    unique (login, session_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema users_schema cascade;
drop schema auth_schema cascade;
-- +goose StatementEnd
