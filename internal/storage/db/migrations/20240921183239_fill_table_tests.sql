-- +goose Up
-- +goose StatementBegin
insert into test(id, name, description, diff_level, category)
values (1, 'Угадай персонажей аниме Наруто',
        'Попробуй угадать персонажа аниме и манги Наруто!',
       1,
        'аниме');
insert into test(id, name, description, diff_level, category)
values (2, 'Угадай персонажей советских мультфильмов',
        'Попробуй угадать персонажей из советских мультфильмов!',
        2,
        'мультфильмы');
insert into test(id, name, description, diff_level, category)
values (3, 'Школьный тест по инфоматике',
        'В тесте собраны вопросы из ОГЭ по информатике, проверьте свои знания!',
        1,
        'школа');


insert into test_questions(id, test_id, question, is_song)
values (1, 1, 'Как зовут первого учителя Наруто?', false);

insert into question_variants(question_id, answer, is_correct)
values (1, 'Джирайя', false);
insert into question_variants(question_id, answer, is_correct)
values (1, 'Саске', false);
insert into question_variants(question_id, answer, is_correct)
values (1, 'Ирука', true);
insert into question_variants(question_id, answer, is_correct)
values (1, 'Конохамару', false);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table test cascade;
truncate table test_questions cascade;
truncate table question_variants cascade;
-- +goose StatementEnd
