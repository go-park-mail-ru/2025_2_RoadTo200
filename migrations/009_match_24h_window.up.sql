-- Миграция 009: Добавление механики 24-часового окна для мэтчей

-- Добавляем поле expires_at для отслеживания времени истечения 24-часового окна
ALTER TABLE match ADD COLUMN expires_at TIMESTAMPTZ;

-- Для существующих мэтчей устанавливаем expires_at как matched_at + 24 часа
-- Если есть хотя бы одно сообщение, устанавливаем NULL (матч активен навсегда)
UPDATE match m
SET expires_at = CASE 
    WHEN EXISTS (SELECT 1 FROM message WHERE match_id = m.id) THEN NULL
    ELSE m.matched_at + INTERVAL '24 hours'
END;

-- Деактивируем мэтчи, у которых истекло время и нет сообщений
UPDATE match
SET is_active = FALSE
WHERE expires_at IS NOT NULL 
  AND expires_at < NOW()
  AND is_active = TRUE;

-- Создаем индекс для быстрого поиска истекших мэтчей
CREATE INDEX IF NOT EXISTS idx_match_expires_at ON match(expires_at) WHERE expires_at IS NOT NULL AND is_active = TRUE;

-- Комментарии
COMMENT ON COLUMN match.expires_at IS 'Время истечения 24-часового окна. NULL если кто-то написал сообщение (матч активен навсегда)';

