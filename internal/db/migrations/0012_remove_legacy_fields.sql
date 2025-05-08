-- Включаем безопасное удаление (без ошибок при отсутствии объектов)
SET client_min_messages TO WARNING;

-- Удаляем устаревшие текстовые поля из таблицы reviews
-- Теперь используем только ссылки по ID на справочные таблицы
ALTER TABLE reviews DROP COLUMN employment;
ALTER TABLE reviews DROP COLUMN employment_period;
ALTER TABLE reviews DROP COLUMN city;

-- Записываем информацию о миграции
INSERT INTO migration_history (filename) VALUES ('0012_remove_legacy_fields.sql'); 