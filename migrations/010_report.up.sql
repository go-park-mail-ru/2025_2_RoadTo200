-- Миграция 010: Создание таблицы отчетов

-- Создаем типы только если они не существуют
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'report_theme_enum') THEN
        CREATE TYPE report_theme_enum AS ENUM ('technical', 'feature', 'question', 'security', 'billing', 'device');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'report_status_enum') THEN
        CREATE TYPE report_status_enum AS ENUM ('open', 'work', 'close');
    END IF;
END $$;

-- Создаем таблицу только если она не существует
CREATE TABLE IF NOT EXISTS report
(
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    theme   report_theme_enum NOT NULL,
    problem TEXT NOT NULL,
    contact TEXT NOT NULL,
    comment TEXT,
    status report_status_enum NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    work_at TIMESTAMPTZ DEFAULT NULL,
    closed_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT report_contact_check CHECK (contact ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'),
    CONSTRAINT report_comment_check CHECK (LENGTH(comment) BETWEEN 1 AND 250),
    CONSTRAINT report_problem_check CHECK (LENGTH(problem) BETWEEN 1 AND 250)
);