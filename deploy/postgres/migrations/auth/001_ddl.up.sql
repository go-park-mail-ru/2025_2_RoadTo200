-- Миграция 001: Создание начальной схемы базы данных

-- Подключаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Создаем пользовательские ENUM типы для повышения целостности данных
CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other');

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
    city        TEXT,
    artist      TEXT,
    quote       TEXT,
    is_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_email_check CHECK (email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'),
    CONSTRAINT user_age_check CHECK (birth_date <= (NOW() - INTERVAL '18 years')::date),
    CONSTRAINT user_name_length_check CHECK (LENGTH(name) BETWEEN 1 AND 50),
    CONSTRAINT user_password_length_check CHECK (LENGTH(TRIM(password)) BETWEEN 8 AND 60),
    CONSTRAINT user_phone_format_check CHECK (phone IS NULL OR phone ~ '^\+?[0-9\s\-\(\)]{10,20}$'),
    CONSTRAINT user_bio_length_check CHECK (LENGTH(TRIM(bio)) < 255),
    CONSTRAINT user_city_check CHECK (LENGTH(city) BETWEEN 1 AND 50),
    CONSTRAINT user_artist_check CHECK (LENGTH(artist) BETWEEN 1 AND 50),
    CONSTRAINT user_quote_check CHECK (LENGTH(quote) BETWEEN 1 AND 50)
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

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE
    ON "user"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_user_email ON "user"(email);
CREATE INDEX IF NOT EXISTS idx_user_phone ON "user"(phone);

SELECT 'DDL executed successfully' as status;