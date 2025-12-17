-- Скрипт генерации тестовых данных

-- Вставка 15 тестовых пользователей
INSERT INTO "user" (email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active) VALUES
-- Мужчины
('alex.ivanov@test.com', '+79161234567', 'Алексей Иванов', 'password123', '1990-05-15', 'male', 'Люблю путешествия и спорт', 'Москва', 'The Weeknd', 'Живи, а работай в свободное время', true, NOW() - INTERVAL '2 hours'),
('max.petrov@test.com', '+79161234568', 'Максим Петров', 'password123', '1988-12-20', 'male', 'Ищу серьезные отношения', 'Санкт-Петербург', 'Arctic Monkeys', 'Всё или ничего', false, NOW() - INTERVAL '1 day'),
('dmitry.sidorov@test.com', '+79161234569', 'Дмитрий Сидоров', 'password123', '1995-08-10', 'male', 'Фотограф, музыкант, мечтатель', 'Казань', 'Tame Impala', 'Красота в простом', true, NOW() - INTERVAL '30 minutes'),
('ivan.kuznetsov@test.com', '+79161234570', 'Иван Кузнецов', 'password123', '1992-03-25', 'male', 'Люблю активный отдых', 'Новосибирск', 'Daft Punk', 'Технологии меняют мир', false, NOW() - INTERVAL '5 hours'),
('sergey.popov@test.com', '+79161234571', 'Сергей Попов', 'password123', '1985-11-08', 'male', 'Бизнесмен, инвестор', 'Екатеринбург', 'Queen', 'Успех - это путь', true, NOW() - INTERVAL '3 days'),

-- Женщины
('anna.smirnova@test.com', '+79161234572', 'Анна Смирнова', 'password123', '1993-07-12', 'female', 'Люблю искусство и кофе', 'Москва', 'Taylor Swift', 'Мечтайте о великом', true, NOW() - INTERVAL '1 hour'),
('maria.volkova@test.com', '+79161234573', 'Мария Волкова', 'password123', '1991-09-18', 'female', 'Психолог, йога-инструктор', 'Санкт-Петербург', 'Lana Del Rey', 'Внутренний мир важнее внешнего', false, NOW() - INTERVAL '2 days'),
('elena.kuzmina@test.com', '+79161234574', 'Елена Кузьмина', 'password123', '1994-04-30', 'female', 'Дизайнер, иллюстратор', 'Краснодар', 'Billie Eilish', 'Творчество - моя жизнь', true, NOW()),
('olga.romanova@test.com', '+79161234575', 'Ольга Романова', 'password123', '1989-01-22', 'female', 'Юрист, люблю театр', 'Ростов-на-Дону', 'Adele', 'Закон должен служить людям', false, NOW() - INTERVAL '8 hours'),
('ekaterina.fedorova@test.com', '+79161234576', 'Екатерина Федорова', 'password123', '1996-12-05', 'female', 'Студентка, модель', 'Сочи', 'Rihanna', 'Уверенность - ключ к успеху', true, NOW() - INTERVAL '45 minutes'),

-- Другие пользователи
('alexandra.novikova@test.com', '+79161234577', 'Александра Новикова', 'password123', '1997-06-14', 'female', 'Фитнес-тренер', 'Москва', 'Beyoncé', 'Здоровье - главное богатство', true, NOW() - INTERVAL '6 hours'),
('viktor.orlov@test.com', '+79161234578', 'Виктор Орлов', 'password123', '1987-02-28', 'male', 'IT-специалист, геймер', 'Санкт-Петербург', 'Hans Zimmer', 'Код - это поэзия', false, NOW() - INTERVAL '1 day'),
('natalia.lebedeva@test.com', '+79161234579', 'Наталья Лебедева', 'password123', '1990-11-11', 'female', 'Врач, волонтер', 'Владивосток', 'Coldplay', 'Помогать людям - мое призвание', true, NOW() - INTERVAL '12 hours'),
('artem.morozov@test.com', '+79161234580', 'Артем Морозов', 'password123', '1993-08-03', 'male', 'Предприниматель', 'Калининград', 'The Beatles', 'Инновации меняют будущее', false, NOW() - INTERVAL '4 hours'),
('irina.pavlova@test.com', '+79161234581', 'Ирина Павлова', 'password123', '1998-04-17', 'female', 'Студентка, художница', 'Самара', 'Lord Huron', 'Искусство вечно', true, NOW() - INTERVAL '20 minutes');

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

-- Вставка подписок
INSERT INTO subscription (user_id, plan_type, start_date, end_date, is_active) VALUES
        ((SELECT id FROM "user" WHERE email = 'alex.ivanov@test.com'), 'premium', NOW(), NOW() + INTERVAL '30 days', true),
        ((SELECT id FROM "user" WHERE email = 'anna.smirnova@test.com'), 'gold', NOW(), NOW() + INTERVAL '60 days', true),
        ((SELECT id FROM "user" WHERE email = 'elena.kuzmina@test.com'), 'platinum', NOW(), NOW() + INTERVAL '90 days', true);

SELECT 'Test data inserted successfully' as status;