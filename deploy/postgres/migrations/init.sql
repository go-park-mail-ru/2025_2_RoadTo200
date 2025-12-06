-- 1. Создаем пользователя для приложения
CREATE USER app_service WITH
    PASSWORD 'sysadm'
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOINHERIT
    NOREPLICATION
    CONNECTION LIMIT 100;

COMMENT ON ROLE app_service IS 'Сервисная учетная запись для основного приложения';

-- 2. Даем права на подключение к базе данных
GRANT CONNECT ON DATABASE "Tinder" TO app_service;

-- 3. Даем права на использование схемы public
GRANT USAGE ON SCHEMA public TO app_service;

-- 8. Настраиваем дефолтные права для будущих таблиц
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_service;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT EXECUTE ON FUNCTIONS TO app_service;

-- 9. Отзыв лишних прав
REVOKE ALL ON DATABASE "Tinder" FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC;
REVOKE ALL ON ALL FUNCTIONS IN SCHEMA public FROM PUBLIC;

CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- 1. Создаем пользователя для мониторинга
CREATE USER monitor_user WITH
    PASSWORD '${MONITOR_PASSWORD}'  -- Пароль из переменной окружения
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOINHERIT
    NOREPLICATION
    CONNECTION LIMIT 10;

COMMENT ON ROLE monitor_user IS 'Пользователь для систем мониторинга (Prometheus, Grafana)';

-- 2. Даем права на подключение к базе данных
GRANT CONNECT ON DATABASE "Tinder" TO monitor_user;

-- 3. Даем роль pg_monitor (включает права на статистику)
GRANT pg_monitor TO monitor_user;

-- 4. Дополнительные права для мониторинга

-- Доступ к схеме public
GRANT USAGE ON SCHEMA public TO monitor_user;

-- Чтение всех таблиц
GRANT SELECT ON ALL TABLES IN SCHEMA public TO monitor_user;

-- 5. Даем доступ к pg_stat_statements
GRANT SELECT ON pg_stat_statements TO monitor_user;

-- 6. Настройка для будущих таблиц
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT ON TABLES TO monitor_user;

-- Проверяем, что у пользователя есть нужные права
SELECT
    table_schema,
    table_name,
    array_agg(privilege_type) as privileges
FROM information_schema.role_table_grants
WHERE grantee = 'app_service'
GROUP BY table_schema, table_name
ORDER BY table_schema, table_name;
