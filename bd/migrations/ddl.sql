-- Миграция 001: Создание начальной схемы базы данных

-- Подключаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Создаем пользовательские ENUM типы для повышения целостности данных
CREATE TYPE gender_enum AS ENUM ('male', 'female');
CREATE TYPE gender_preference_enum AS ENUM ('male', 'female');
CREATE TYPE swipe_type_enum AS ENUM ('like', 'dislike', 'super_like');
CREATE TYPE plan_type_enum AS ENUM ('premium', 'gold', 'platinum');

-- Таблица: user
-- Используем TEXT вместо VARCHAR(n) и TIMESTAMPTZ для корректной работы с часовыми поясами.
CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE,
    name TEXT NOT NULL,
    birth_date DATE NOT NULL,
    gender gender_enum NOT NULL,
    bio TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_email_check CHECK (email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'),
    CONSTRAINT user_age_check CHECK (birth_date <= (NOW() - INTERVAL '18 years')::date)
);

-- Таблица: user_photo
CREATE TABLE user_photo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    display_order SMALLINT NOT NULL DEFAULT 0,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON user_photo (user_id);

-- Таблица: user_preference
-- user_id является одновременно и первичным, и внешним ключом для связи 1-к-1
CREATE TABLE user_preference (
    user_id UUID PRIMARY KEY REFERENCES "user"(id) ON DELETE CASCADE,
    show_gender gender_preference_enum NOT NULL DEFAULT 'both',
    age_min SMALLINT NOT NULL DEFAULT 18,
    age_max SMALLINT NOT NULL DEFAULT 99,
    max_distance INT NOT NULL DEFAULT 100,
    global_search BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_preference_age_range_check CHECK (age_min >= 18 AND age_max > age_min AND age_max <= 120),
    CONSTRAINT user_preference_distance_check CHECK (max_distance > 0)
);

-- Таблица: swipe
-- Составной первичный ключ гарантирует, что один пользователь может свайпнуть другого только один раз.
CREATE TABLE swipe (
    swiper_user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    target_user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    swipe_type swipe_type_enum NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (swiper_user_id, target_user_id),
    CONSTRAINT swipe_no_self_swipe_check CHECK (swiper_user_id <> target_user_id)
);

-- Индекс для быстрого поиска тех, кто свайпнул пользователя
CREATE INDEX ON swipe (target_user_id);

-- Таблица: match
CREATE TABLE match (
    user1_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user1_id, user2_id),
    CONSTRAINT match_no_self_match_check CHECK (user1_id <> user2_id)
);

CREATE INDEX ON match (user1_id);
CREATE INDEX ON match (user2_id);

-- Таблица: message
CREATE TABLE message (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    message_text TEXT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT message_text_length_check CHECK (char_length(message_text) > 0 AND char_length(message_text) <= 1000),
    CONSTRAINT user_preference_age_range_check CHECK (0 <= status AND status <= 2),
);

CREATE INDEX ON message (match_id, created_at DESC);

-- Таблица: subscription
CREATE TABLE subscription (
    user_id UUID PRIMARY KEY REFERENCES "user"(id) ON DELETE CASCADE,
    plan_type plan_type_enum NOT NULL,
    start_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_date TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT subscription_date_check CHECK (end_date > start_date)
);
