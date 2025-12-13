-- Миграция 005: Обновление enum планов подписки

-- Сначала изменяем колонку на text, чтобы можно было изменить enum
-- Если есть старые данные, удаляем их (для MVP можно удалить все старые подписки)
DELETE FROM subscription WHERE plan_type IN ('premium', 'gold', 'platinum');

ALTER TABLE subscription ALTER COLUMN plan_type TYPE text;

-- Удаляем старый enum (если есть зависимые объекты, они уже удалены)
DROP TYPE IF EXISTS plan_type_enum;

-- Создаем новый enum с актуальными значениями
CREATE TYPE plan_type_enum AS ENUM ('week', 'month', 'quarter');

-- Возвращаем колонку к новому enum типу
-- Если остались какие-то значения, они будут конвертированы, но лучше чтобы их не было
ALTER TABLE subscription ALTER COLUMN plan_type TYPE plan_type_enum USING plan_type::plan_type_enum;

