-- Миграция 008: Добавление типа уведомления 'message'

ALTER TYPE notification_type_enum ADD VALUE 'message';

