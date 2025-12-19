-- Откат миграции 010: Удаление таблицы отчетов

DROP TABLE IF EXISTS report;
DROP TYPE IF EXISTS report_theme_enum;
DROP TYPE IF EXISTS report_status_enum;;

