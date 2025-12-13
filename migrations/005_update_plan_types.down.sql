-- Откат миграции 005: Возврат старых значений enum планов

ALTER TABLE subscription ALTER COLUMN plan_type TYPE text;

DROP TYPE plan_type_enum;

CREATE TYPE plan_type_enum AS ENUM ('premium', 'gold', 'platinum');

ALTER TABLE subscription ALTER COLUMN plan_type TYPE plan_type_enum USING plan_type::plan_type_enum;

