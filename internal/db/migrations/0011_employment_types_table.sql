-- Включаем безопасное удаление (без ошибок при отсутствии объектов)
SET client_min_messages TO WARNING;

-- Создаем таблицу для хранения типов занятости
CREATE TABLE IF NOT EXISTS employment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

-- Добавляем базовые типы занятости
INSERT INTO employment_types (name, description) VALUES 
    ('Удалённый формат работы', 'Работа выполняется вне офиса компании'),
    ('Офисный формат', 'Работа выполняется в офисе компании'),
    ('Полная занятость', 'Работа на полный рабочий день'),
    ('Частичная занятость', 'Работа на неполный рабочий день'),
    ('Гибридный формат', 'Сочетание удалённой и офисной работы'),
    ('Проектная работа', 'Работа по проектам с ограниченным сроком'),
    ('Фриланс', 'Независимая работа без постоянного трудоустройства')
ON CONFLICT (name) DO NOTHING;

-- Создаем временную таблицу для хранения существующих значений типов занятости
CREATE TEMPORARY TABLE temp_reviews_employment AS
SELECT id, employment FROM reviews;

-- Добавляем колонку employment_type_id в таблицу reviews
ALTER TABLE reviews ADD COLUMN employment_type_id INTEGER REFERENCES employment_types(id);

-- Обновляем данные в таблице reviews
UPDATE reviews r SET employment_type_id = (
    SELECT et.id FROM employment_types et 
    WHERE temp.employment ILIKE '%' || et.name || '%' OR et.name ILIKE '%' || temp.employment || '%'
    LIMIT 1
)
FROM temp_reviews_employment temp
WHERE r.id = temp.id;

-- Удаляем временную таблицу
DROP TABLE temp_reviews_employment;

-- Записываем информацию о миграции
INSERT INTO migration_history (filename) VALUES ('0011_employment_types_table.sql'); 