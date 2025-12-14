-- Миграция 007: Создание таблицы уведомлений

CREATE TYPE notification_type_enum AS ENUM ('match', 'super_like', 'like');

CREATE TABLE notification (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    type       notification_type_enum NOT NULL,
    from_user_id UUID REFERENCES "user" (id) ON DELETE CASCADE,
    match_id   UUID REFERENCES match (id) ON DELETE CASCADE,
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX idx_notification_user_id ON notification(user_id);
CREATE INDEX idx_notification_user_unread ON notification(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notification_created_at ON notification(created_at DESC);

