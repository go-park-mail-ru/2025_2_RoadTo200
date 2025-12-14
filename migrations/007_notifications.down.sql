-- Откат миграции 007: Удаление таблицы уведомлений

DROP TABLE IF EXISTS notification;
DROP TYPE IF EXISTS notification_type_enum;

