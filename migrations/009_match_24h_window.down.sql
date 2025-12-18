-- Откат миграции 009: Удаление механики 24-часового окна

-- Удаляем индекс
DROP INDEX IF EXISTS idx_match_expires_at;

-- Удаляем поле expires_at
ALTER TABLE match DROP COLUMN IF EXISTS expires_at;

-- Восстанавливаем все мэтчи как активные (опционально)
-- UPDATE match SET is_active = TRUE WHERE NOT is_active;

