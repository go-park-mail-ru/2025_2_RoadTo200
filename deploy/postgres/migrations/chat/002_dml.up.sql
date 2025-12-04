-- Вставка мэтчей (когда оба пользователя лайкнули друг друга)
INSERT INTO match (user1_id, user2_id, is_active) VALUES
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), (SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), true),
        ((SELECT id FROM "user" WHERE email = 'max.petrov@test.com'), (SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'), true);

-- Вставка сообщений в мэтчи (обновлено для новой схемы)
INSERT INTO message (match_id, sender_id, receiver_id, content, is_read) VALUES
        ((SELECT id FROM match WHERE user1_id = (SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com') AND user2_id = (SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com')), 
         (SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 
         (SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'),
         'Привет! Как дела?', false),
        ((SELECT id FROM match WHERE user1_id = (SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com') AND user2_id = (SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com')), 
         (SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'),
         (SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'),
         'Привет! Все отлично, а у тебя?', false),
        ((SELECT id FROM match WHERE user1_id = (SELECT id FROM "user" WHERE email = 'max.petrov@test.com') AND user2_id = (SELECT id FROM "user" WHERE email = 'maria.volkova@test.com')), 
         (SELECT id FROM "user" WHERE email = 'max.petrov@test.com'),
         (SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'),
         'Очень рад нашему знакомству!', false);

SELECT 'Test data inserted successfully' as status;