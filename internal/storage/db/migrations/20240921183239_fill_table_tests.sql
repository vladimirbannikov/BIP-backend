-- +goose Up
-- +goose StatementBegin

/*tests*/

insert into test(id, name, description, diff_level, category, pictureFile)
values (1, 'Угадай персонажей аниме Наруто',
        'Попробуй угадать персонажа аниме и манги Наруто!',
       1,
        'аниме',
        'naruto-test.jpg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (2, 'Угадай персонажей советских мультфильмов',
        'Попробуй угадать персонажей из советских мультфильмов!',
        2,
        'мультфильмы',
        'soviet-mult.jpg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (3, 'Школьный тест по инфоматике',
        'В тесте собраны вопросы из ОГЭ по информатике, проверьте свои знания!',
        1,
        'школа',
        'info-oge.jpg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (4, 'Угадай персонажей американских мультфильмов',
        'Попробуй угадать персонажей из американских мультфильмов!',
        2,
        'мультфильмы',
        'american_shit.jpg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (5, 'Школьный тест по биологии',
        'В тесте собраны вопросы из ОГЭ по биологии, проверьте свои знания!',
        3,
        'школа',
        'biotest.jpg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (6, 'anime test 2',
        'anime test  2 description',
        1,
        'аниме',
        'animetest2.jpeg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (7, 'anime test 3',
        'anime test  3 description',
        1,
        'аниме',
        'animetest3.jpeg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (8, 'anime test 4',
        'anime test 4 description',
        1,
        'аниме',
        'animetest4.jpeg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (9, 'anime test 5',
        'anime test  5 description',
        1,
        'аниме',
        'animetest5.jpeg');
insert into test(id, name, description, diff_level, category, pictureFile)
values (10, 'Школьный тест по математике',
        'В тесте собраны вопросы из ОГЭ по математике, проверьте свои знания!',
    1,
        'школа',
        'mathtest.jpg');

insert into test(id, name, description, diff_level, category, pictureFile)
values (11, 'Тест на знатока музыки восьмедисятых',
        'В тесте собраны интересные вопросы на тему музыки восьмедисятых',
        3,
        'музыка',
        'music_test.jpg');


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


insert into test_questions(id, test_id, question, is_song, song_file)
values (5, 11, 'Угадай песню:', true, 'messagesfromthestars.mp3');

insert into question_variants(question_id, answer, is_correct)
values (5, 'A) Джастин Бибер - Бейби', false);
insert into question_variants(question_id, answer, is_correct)
values (5, 'Б) Нюша - Ночь', false);
insert into question_variants(question_id, answer, is_correct)
values (5, 'В) The RAH Band - Messages From the Stars', true);
insert into question_variants(question_id, answer, is_correct)
values (5, 'Г) Михаил Круг - Золотые купола', false);

insert into test_questions(id, test_id, question, is_song, song_file)
values (6, 11, 'Угадай песню:', true, 'smoothoperator.mp3');

insert into question_variants(question_id, answer, is_correct)
values (6, 'A) Авария - Пей пиво', false);
insert into question_variants(question_id, answer, is_correct)
values (6, 'Б) Sade - Smooth Operator', true);
insert into question_variants(question_id, answer, is_correct)
values (6, 'В) 2pac - Changez', false);
insert into question_variants(question_id, answer, is_correct)
values (6, 'Г) Руки вверх! — Он тебя целует', false);

insert into test_questions(id, test_id, question, is_song, song_file)
values (7, 11, 'Угадай песню:', true, 'funkytown.mp3');

insert into question_variants(question_id, answer, is_correct)
values (7, 'A) Lipps Inc - Funky town', true);
insert into question_variants(question_id, answer, is_correct)
values (7, 'Б) Вирус - Ручки', false);
insert into question_variants(question_id, answer, is_correct)
values (7, 'В) A-ha - Take on me', false);
insert into question_variants(question_id, answer, is_correct)
values (7, 'Г) Бригада - Интро', false);

insert into test_questions(id, test_id, question, is_song)
values (8, 11, 'В каком году распалась группа Битлз?', false);

insert into question_variants(question_id, answer, is_correct)
values (8, '1998', false);
insert into question_variants(question_id, answer, is_correct)
values (8, '1970', true);
insert into question_variants(question_id, answer, is_correct)
values (8, '1853', false);
insert into question_variants(question_id, answer, is_correct)
values (8, '2018', false);

/*users, score*/
insert into users_schema.users(login, password_hash, email)
values ('user1', 'passhash', 'user1@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score, avatarFile)
values ('user1', 2, 5, 'file_149450541.jpg');
insert into user_test_score(test_id, user_login, score)
values (1, 'user1', 1);
insert into user_test_score(test_id, user_login, score)
values (2, 'user1', 4);

insert into users_schema.users(login, password_hash, email)
values ('user2', 'passhash', 'user2@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score, avatarFile)
values ('user2', 1, 6, 'file_149450563.jpg');
insert into user_test_score(test_id, user_login, score)
values (2, 'user2',6);

insert into users_schema.users(login, password_hash, email)
values ('user3', 'passhash', 'user3@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score, avatarFile)
values ('user3', 0, 0, 'file_149450558.jpg');

insert into users_schema.users(login, password_hash, email)
values ('user4', 'passhash', 'user4@gmail.com');
insert into users_schema.user_profile(login, tests_count, total_score, avatarFile)
values ('user4', 2, 0, 'file_149450546.jpg');
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
