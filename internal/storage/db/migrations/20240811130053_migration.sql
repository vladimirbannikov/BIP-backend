-- +goose Up
-- +goose StatementBegin
alter TABLE users_schema.user_profile
      add column tests_count int;
alter TABLE users_schema.user_profile
      add column total_score int;
/*alter table users_schema.user_profile
    add column avatar text;*/

create table if not exists test (
    id serial primary key,
    name text unique,
    description text,
    diff_level int,
    category text
);

create table if not exists test_questions (
    id serial primary key,
    test_id serial REFERENCES test(id),
    question text,
    is_song bool
);

create table if not exists question_variants (
    id serial primary key,
    question_id serial REFERENCES test_questions(id),
    answer text,
    is_correct bool
);

create table if not exists user_test_score
(
    id         serial primary key,
    test_id    serial REFERENCES test (id),
    user_login text REFERENCES users_schema.users (login),
    score int,
    CONSTRAINT unique_user_test_score unique(test_id, user_login)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users_schema.user_profile
    DROP COLUMN tests_count;
ALTER TABLE users_schema.user_profile
    DROP COLUMN total_score;
/*ALTER TABLE users_schema.user_profile
    DROP COLUMN avatar;*/

drop table test cascade ;
drop table test_questions cascade;
drop table question_variants cascade;
drop table user_test_score cascade;

-- +goose StatementEnd
