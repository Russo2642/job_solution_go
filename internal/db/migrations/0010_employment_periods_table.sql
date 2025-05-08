-- Включаем безопасное удаление (без ошибок при отсутствии объектов)
SET client_min_messages TO WARNING;

-- Создаем таблицу для хранения периодов работы
CREATE TABLE IF NOT EXISTS employment_periods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

-- Добавляем базовые периоды работы
INSERT INTO employment_periods (name, description) VALUES 
    ('Меньше года', 'Опыт работы в компании менее 1 года'),
    ('1-3 года', 'Опыт работы в компании от 1 до 3 лет'),
    ('3-5 лет', 'Опыт работы в компании от 3 до 5 лет'),
    ('Больше 5 лет', 'Опыт работы в компании более 5 лет')
ON CONFLICT (name) DO NOTHING;

-- Создаем временную таблицу для хранения существующих значений периодов работы
CREATE TEMPORARY TABLE temp_reviews_employment_period AS
SELECT id, employment_period FROM reviews;

-- Добавляем колонку employment_period_id в таблицу reviews
ALTER TABLE reviews ADD COLUMN employment_period_id INTEGER REFERENCES employment_periods(id);

-- Обновляем данные в таблице reviews
UPDATE reviews r SET employment_period_id = (
    SELECT ep.id FROM employment_periods ep 
    WHERE temp.employment_period ILIKE '%' || ep.name || '%' OR ep.name ILIKE '%' || temp.employment_period || '%'
    LIMIT 1
)
FROM temp_reviews_employment_period temp
WHERE r.id = temp.id;

-- Удаляем временную таблицу
DROP TABLE temp_reviews_employment_period;

-- Записываем информацию о миграции
INSERT INTO migration_history (filename) VALUES ('0010_employment_periods_table.sql'); 