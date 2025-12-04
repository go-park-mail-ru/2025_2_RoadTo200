-- Вставка фотографий пользователей
INSERT INTO user_photo (user_id, photo_url, display_order, is_approved) VALUES
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'https://example.com/photos/alex1.jpg', 1, true),
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'https://example.com/photos/alex2.jpg', 2, true),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'https://example.com/photos/anna1.jpg', 1, true),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'https://example.com/photos/anna2.jpg', 2, true),
        ((SELECT id FROM "user" WHERE email = 'max.petrov@test.com'), 'https://example.com/photos/max1.jpg', 1, true),
        ((SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'), 'https://example.com/photos/maria1.jpg', 1, true),
        ((SELECT id FROM "user" WHERE email = 'dmitry.sidorov@test.com'), 'https://example.com/photos/dmitry1.jpg', 1, true),
        ((SELECT id FROM "user" WHERE email = 'elena.kuzmina@test.com'), 'https://example.com/photos/elena1.jpg', 1, true);

-- Вставка предпочтений пользователей
INSERT INTO user_preference (user_id, show_gender, age_min, age_max, max_distance, global_search) VALUES
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'female', 25, 35, 50, false),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'male', 28, 40, 30, false),
        ((SELECT id FROM "user" WHERE email = 'max.petrov@test.com'), 'female', 23, 33, 100, true),
        ((SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'), 'male', 30, 45, 25, false),
        ((SELECT id FROM "user" WHERE email = 'dmitry.sidorov@test.com'), 'both', 20, 30, 75, false);

-- Вставка интересов
INSERT INTO interest (user_id, theme) VALUES
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'workout'),
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'chill'),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'culture'),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'fun'),
        ((SELECT id FROM "user" WHERE email = 'max.petrov@test.com'), 'relax'),
        ((SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'), 'yoga'),
        ((SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'), 'love'),
        ((SELECT id FROM "user" WHERE email = 'dmitry.sidorov@test.com'), 'cinema'),
        ((SELECT id FROM "user" WHERE email = 'elena.kuzmina@test.com'), 'culture'),
        ((SELECT id FROM "user" WHERE email = 'elena.kuzmina@test.com'), 'cinema');

-- Вставка свайпов
INSERT INTO swipe (swiper_user_id, target_user_id, swipe_type) VALUES
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), (SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'like'),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), (SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'like'),
        ((SELECT id FROM "user" WHERE email = 'max.petrov@test.com'), (SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'), 'like'),
        ((SELECT id FROM "user" WHERE email = 'dmitry.sidorov@test.com'), (SELECT id FROM "user" WHERE email = 'elena.kuzmina@test.com'), 'super_like'),
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), (SELECT id FROM "user" WHERE email = 'maria.volkova@test.com'), 'dislike');

-- Вставка подписок
INSERT INTO subscription (user_id, plan_type, start_date, end_date, is_active) VALUES
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'premium', NOW(), NOW() + INTERVAL '30 days', true),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'gold', NOW(), NOW() + INTERVAL '60 days', true),
        ((SELECT id FROM "user" WHERE email = 'elena.kuzmina@test.com'), 'platinum', NOW(), NOW() + INTERVAL '90 days', true);

SELECT 'Test data inserted successfully' as status;