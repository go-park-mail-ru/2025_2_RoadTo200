-- Откат миграции 004: Удаление полей премиум-функционала

ALTER TABLE "user" DROP COLUMN IF EXISTS super_likes_count;
ALTER TABLE "user" DROP COLUMN IF EXISTS is_premium;

