-- Миграция 001: Создание начальной схемы базы данных

-- Подключаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Создаем пользовательские ENUM типы для повышения целостности данных
CREATE TYPE gender_preference_enum AS ENUM ('male', 'female', 'both');
CREATE TYPE swipe_type_enum AS ENUM ('like', 'dislike', 'super_like');
CREATE TYPE plan_type_enum AS ENUM ('premium', 'gold', 'platinum');
CREATE TYPE interest_theme_enum AS ENUM ('workout', 'fun', 'party', 'chill', 'love', 'relax', 'yoga', 'friendship', 'culture', 'cinema');
CREATE TYPE strike_reason_type AS ENUM ('spam', 'fake_profile', 'offensive_content', 'harassment', 'inappropriate_content', 'underage', 'copyright_violation', 'other');
CREATE TYPE strike_status_type AS ENUM ('pending', 'approved', 'rejected', 'resolved');

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

CREATE TABLE interest
(
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    theme   interest_theme_enum NOT NULL,

    UNIQUE (user_id, theme)
);

CREATE TABLE strike (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        reporter_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
                        target_user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
                        type strike_reason_type NOT NULL,
                        reason TEXT,
                        status strike_status_type NOT NULL DEFAULT 'pending',
                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        updated_at TIMESTAMPTZ,
                        moderator_id UUID REFERENCES "user" (id) ON DELETE CASCADE,
                        moderator_note TEXT,

                        UNIQUE (reporter_id, target_user_id),
                        CONSTRAINT strike_reason_check CHECK (LENGTH(reason) BETWEEN 1 AND 250),
                        CONSTRAINT strike_note_check CHECK (LENGTH(moderator_note) BETWEEN 1 AND 250),
                        CONSTRAINT chk_strikes_dates CHECK (created_at <= COALESCE(updated_at, NOW())),
                        CONSTRAINT chk_strikes_self_report CHECK (reporter_id != target_user_id)
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
-- end;

CREATE TRIGGER update_strike_updated_at
    BEFORE UPDATE
    ON strike
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_user_photo_user_id ON user_photo(user_id);
CREATE INDEX IF NOT EXISTS idx_swipe_swiper_id ON swipe(swiper_user_id);
CREATE INDEX IF NOT EXISTS idx_swipe_target_id ON swipe(target_user_id);
CREATE INDEX idx_message_match_id_created_at ON message(match_id, created_at DESC);
CREATE INDEX idx_message_receiver_unread ON message(receiver_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_message_created_at ON message(created_at DESC);

SELECT 'DDL executed successfully' as status;