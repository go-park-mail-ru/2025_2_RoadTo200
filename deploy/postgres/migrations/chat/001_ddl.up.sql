-- Миграция 001: Создание начальной схемы базы данных

-- Подключаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Создаем пользовательские ENUM типы для повышения целостности данных
CREATE TYPE swipe_type_enum AS ENUM ('like', 'dislike', 'super_like');

-- Таблица: match
CREATE TABLE match
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id   UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    user2_id   UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT match_unique UNIQUE (user1_id, user2_id),
    CONSTRAINT match_no_self_match_check CHECK (user1_id <> user2_id)
);

-- Таблица: message
CREATE TABLE message
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id    UUID        NOT NULL REFERENCES match (id) ON DELETE CASCADE,
    sender_id   UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    receiver_id UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    content     TEXT        NOT NULL,
    is_read     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_swipe_swiper_id ON swipe(swiper_user_id);
CREATE INDEX IF NOT EXISTS idx_swipe_target_id ON swipe(target_user_id);
CREATE INDEX IF NOT EXISTS idx_message_sender_id ON message(sender_id);
CREATE INDEX IF NOT EXISTS idx_message_match_id ON message(match_id);

SELECT 'DDL executed successfully' as status;