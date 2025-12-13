-- Миграция 004: Добавление полей премиум-функционала

ALTER TABLE "user" ADD COLUMN is_premium BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE "user" ADD COLUMN super_likes_count INT NOT NULL DEFAULT 0;

