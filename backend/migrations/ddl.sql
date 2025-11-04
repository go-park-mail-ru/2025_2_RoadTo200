-- Миграция 001: Создание начальной схемы базы данных

-- Подключаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Создаем пользовательские ENUM типы для повышения целостности данных
CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other');
CREATE TYPE gender_preference_enum AS ENUM ('male', 'female', 'both');
CREATE TYPE swipe_type_enum AS ENUM ('like', 'dislike', 'super_like');
CREATE TYPE plan_type_enum AS ENUM ('premium', 'gold', 'platinum');

-- Таблица: user
CREATE TABLE "user"
(
    id          UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    email       TEXT UNIQUE              NOT NULL,
    phone       TEXT UNIQUE,
    name        TEXT                     NOT NULL,
    password    TEXT                     NOT NULL,
    birth_date  DATE,
    gender      gender_enum,
    bio         TEXT,
    latitude    DECIMAL(10, 8),
    longitude   DECIMAL(11, 8),
    is_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_email_check CHECK (email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'),
    CONSTRAINT user_age_check CHECK (birth_date <= (NOW() - INTERVAL '18 years')::date),
    CONSTRAINT user_name_length_check CHECK (LENGTH(TRIM(name)) BETWEEN 1 AND 50),
    CONSTRAINT user_password_length_check CHECK (LENGTH(TRIM(password)) BETWEEN 8 AND 30),
    CONSTRAINT user_phone_format_check CHECK (phone IS NULL OR phone ~ '^\+?[0-9\s\-\(\)]{10,20}$'),
    CONSTRAINT user_bio_length_check CHECK (LENGTH(TRIM(bio)) < 255),
    CONSTRAINT user_latitude_check CHECK (latitude IS NULL OR (latitude BETWEEN -90 AND 90)),
    CONSTRAINT user_longitude_check CHECK (longitude IS NULL OR (longitude BETWEEN -180 AND 180))
);

-- Таблица: user_photo
CREATE TABLE user_photo
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    photo_url     TEXT        NOT NULL,
    display_order SMALLINT    NOT NULL DEFAULT 0,
    is_approved   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Таблица: user_preference
CREATE TABLE user_preference
(
    user_id       UUID PRIMARY KEY REFERENCES "user" (id) ON DELETE CASCADE,
    show_gender   gender_preference_enum NOT NULL DEFAULT 'both',
    age_min       SMALLINT               NOT NULL DEFAULT 18,
    age_max       SMALLINT               NOT NULL DEFAULT 99,
    max_distance  INT                    NOT NULL DEFAULT 100,
    global_search BOOLEAN                NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ            NOT NULL DEFAULT NOW()
);

-- Таблица: swipe
CREATE TABLE swipe
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    swiper_user_id  UUID            NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    target_user_id  UUID            NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    swipe_type      swipe_type_enum NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    CONSTRAINT swipe_no_self_swipe_check CHECK (swiper_user_id <> target_user_id)
);

-- Таблица: match
CREATE TABLE match
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id   UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    user2_id   UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT match_no_self_match_check CHECK (user1_id <> user2_id)
);

-- Таблица: message
CREATE TABLE message
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id     UUID        NOT NULL REFERENCES match (id) ON DELETE CASCADE,
    sender_id    UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    message_text TEXT        NOT NULL,
    status       SMALLINT    NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Таблица: subscription
CREATE TABLE subscription
(
    user_id    UUID PRIMARY KEY REFERENCES "user" (id) ON DELETE CASCADE,
    plan_type  plan_type_enum NOT NULL,
    start_date TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    end_date   TIMESTAMPTZ    NOT NULL,
    is_active  BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE
    ON "user"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_preference_updated_at
    BEFORE UPDATE
    ON user_preference
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_user_email ON "user"(email);
CREATE INDEX IF NOT EXISTS idx_user_phone ON "user"(phone);
CREATE INDEX IF NOT EXISTS idx_user_photo_user_id ON user_photo(user_id);
CREATE INDEX IF NOT EXISTS idx_swipe_swiper_id ON swipe(swiper_user_id);
CREATE INDEX IF NOT EXISTS idx_swipe_target_id ON swipe(target_user_id);
CREATE INDEX IF NOT EXISTS idx_message_sender_id ON message(sender_id);
CREATE INDEX IF NOT EXISTS idx_message_match_id ON message(match_id);
CREATE INDEX IF NOT EXISTS idx_match_user1_id ON match(user1_id);
CREATE INDEX IF NOT EXISTS idx_match_user2_id ON match(user2_id);

SELECT 'DDL executed successfully' as status;