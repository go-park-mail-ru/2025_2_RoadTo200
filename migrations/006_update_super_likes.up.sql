-- Миграция 006: Обновление суперлайков для существующих пользователей

-- Устанавливаем 3 суперлайка всем пользователям, у которых сейчас 0 (или NULL)
UPDATE "user" 
SET super_likes_count = 3 
WHERE super_likes_count = 0 OR super_likes_count IS NULL;

