-- Миграция 002: Наполнение базы данных начальными данными

-- Создание пользователей
INSERT INTO "user" (email, phone, name, birth_date, gender) VALUES
        ('ivan.petrov@example.com', '+79161234567', 'Иван', '1993-03-22', 'male'),
        ('maria.sidorova@example.com', '+79267654321', 'Мария', '1996-07-11', 'female'),
        ('alexey.smirnov@example.com', '+79035558899', 'Алексей', '1990-12-01', 'male');

-- Создание их предпочтений
INSERT INTO user_preference (user_id, show_gender, age_min, age_max) VALUES
        ((SELECT id FROM "user" WHERE email = 'ivan.petrov@example.com'), 'female', 22, 32),
        ((SELECT id FROM "user" WHERE email = 'maria.sidorova@example.com'), 'male', 25, 35),
        ((SELECT id FROM "user" WHERE email = 'alexey.smirnov@example.com'), 'female', 28, 38);

-- Создание свайпов (Иван лайкнул Марию, Мария лайкнула Ивана)
INSERT INTO swipe (swiper_user_id, target_user_id, swipe_type) VALUES
        ((SELECT id FROM "user" WHERE email = 'ivan.petrov@example.com'), (SELECT id FROM "user" WHERE email = 'maria.sidorova@example.com'), 'like'),
        ((SELECT id FROM "user" WHERE email = 'maria.sidorova@example.com'), (SELECT id FROM "user" WHERE email = 'ivan.petrov@example.com'), 'like'),
        ((SELECT id FROM "user" WHERE email = 'alexey.smirnov@example.com'), (SELECT id FROM "user" WHERE email = 'maria.sidorova@example.com'), 'like');

-- Создание мэтча между Иваном и Марией
INSERT INTO match (user1_id, user2_id) VALUES
    ((SELECT id FROM "user" WHERE email = 'ivan.petrov@example.com'), (SELECT id FROM "user" WHERE email = 'maria.sidorova@example.com'));

-- Добавление сообщений в их мэтч
WITH chat_data AS (
    SELECT
        m.id as match_id,
        u1.id as user1_id,
        u2.id as user2_id
    FROM match m
             JOIN "user" u1 ON u1.id = m.user1_id AND u1.email = 'ivan.petrov@example.com'
             JOIN "user" u2 ON u2.id = m.user2_id AND u2.email = 'maria.sidorova@example.com'
)
INSERT INTO message (match_id, sender_id, message_text) VALUES
        ((SELECT match_id FROM chat_data), (SELECT user1_id FROM chat_data), 'Привет, Мария! Рад нашему мэтчу.'),
        ((SELECT match_id FROM chat_data), (SELECT user2_id FROM chat_data), 'Иван, привет! Взаимно :)');

-- Добавление подписки для Алексея
INSERT INTO subscription (user_id, plan_type, end_date) VALUES
        ((SELECT id FROM "user" WHERE email = 'alexey.smirnov@example.com'), 'gold', NOW() + INTERVAL '1 month');
