-- +goose Up
-- +goose StatementBegin

/*tests*/

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
insert into test(id, name, description, diff_level, category)
values (4, 'Угадай персонажей американских мультфильмов',
        'Попробуй угадать персонажей из американских мультфильмов!',
        2,
        'мультфильмы');
insert into test(id, name, description, diff_level, category)
values (5, 'Школьный тест по биологии',
        'В тесте собраны вопросы из ОГЭ по биологии, проверьте свои знания!',
        3,
        'школа');
insert into test(id, name, description, diff_level, category)
values (6, 'anime test 2',
        'anime test  2 description',
        1,
        'аниме');
insert into test(id, name, description, diff_level, category)
values (7, 'anime test 3',
        'anime test  3 description',
        1,
        'аниме');
insert into test(id, name, description, diff_level, category)
values (8, 'anime test 4',
        'anime test 4 description',
        1,
        'аниме');
insert into test(id, name, description, diff_level, category)
values (9, 'anime test 5',
        'anime test  5 description',
        1,
        'аниме');
insert into test(id, name, description, diff_level, category)
values (10, 'Школьный тест по математике',
        'В тесте собраны вопросы из ОГЭ по математике, проверьте свои знания!',
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


insert into test_questions(id, test_id, question, is_song)
values (2, 6, 'Question 1?', false);

insert into question_variants(question_id, answer, is_correct)
values (2, 'var1', false);
insert into question_variants(question_id, answer, is_correct)
values (2, 'var2', false);
insert into question_variants(question_id, answer, is_correct)
values (2, 'var3', true);
insert into question_variants(question_id, answer, is_correct)
values (2, 'var4', false);

insert into test_questions(id, test_id, question, is_song)
values (3, 6, 'Question 2?', false);

insert into question_variants(question_id, answer, is_correct)
values (3, 'var1', false);
insert into question_variants(question_id, answer, is_correct)
values (3, 'var2', false);
insert into question_variants(question_id, answer, is_correct)
values (3, 'var3', true);
insert into question_variants(question_id, answer, is_correct)
values (3, 'var4', false);

insert into test_questions(id, test_id, question, is_song)
values (4, 6, 'Question 3?', false);

insert into question_variants(question_id, answer, is_correct)
values (4, 'var1', false);
insert into question_variants(question_id, answer, is_correct)
values (4, 'var2', false);
insert into question_variants(question_id, answer, is_correct)
values (4, 'var3', true);
insert into question_variants(question_id, answer, is_correct)
values (4, 'var4', false);


/*users, score*/
insert into users_schema.users(login, password_hash, email)
values ('user1', 'passhash', 'user1@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score)
values ('user1', 2, 5);
insert into user_test_score(test_id, user_login, score)
values (1, 'user1', 1);
insert into user_test_score(test_id, user_login, score)
values (2, 'user1', 4);

insert into users_schema.users(login, password_hash, email)
values ('user2', 'passhash', 'user2@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score)
values ('user2', 1, 6);
insert into user_test_score(test_id, user_login, score)
values (2, 'user2',6);

insert into users_schema.users(login, password_hash, email)
values ('user3', 'passhash', 'user3@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score)
values ('user3', 0, 0);

insert into users_schema.users(login, password_hash, email)
values ('user4', 'passhash', 'user4@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score)
values ('user4', 2, 0);
insert into user_test_score(test_id, user_login, score)
values (1, 'user4', 0);
insert into user_test_score(test_id, user_login, score)
values (2, 'user4', 0);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table test cascade;
truncate table test_questions cascade;
truncate table question_variants cascade;
-- +goose StatementEnd
