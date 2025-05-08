-- Включаем безопасное удаление (без ошибок при отсутствии объектов)
SET client_min_messages TO WARNING;

-- Создаем таблицу для хранения типов бенефитов
CREATE TABLE IF NOT EXISTS benefit_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

-- Добавляем базовые типы бенефитов
INSERT INTO benefit_types (name, description) VALUES 
    ('Медицинская страховка', 'Включая ДМС и страхование для поездок за границу'),
    ('Гибкий график', 'Возможность самостоятельно планировать рабочее время'),
    ('Удаленная работа', 'Возможность работать из дома или другого места'),
    ('Бонусы и премии', 'Дополнительные денежные вознаграждения'),
    ('Корпоративные мероприятия', 'Тимбилдинги, корпоративы, праздники'),
    ('Обучение и развитие', 'Курсы, тренинги, конференции'),
    ('Корпоративный транспорт', 'Развозка сотрудников'),
    ('Питание', 'Бесплатное или субсидированное питание в офисе'),
    ('Спорт', 'Фитнес, спортивные мероприятия'),
    ('Дополнительный отпуск', 'Отпуск сверх предусмотренного законодательством'),
    ('Материальная помощь', 'Финансовая поддержка в сложных ситуациях'),
    ('Оплата переезда', 'Оплата переезда в другую местность'),
    ('Оплата обучения', 'Оплата обучения в университете или школе'),
    ('Своевременная оплата труда', 'Оплата труда в установленные сроки'),
    ('Современный офис', 'Офис с современным оборудованием и условиями')
ON CONFLICT (name) DO NOTHING;

-- Создаем временную таблицу для хранения существующих бенефитов отзывов
CREATE TEMPORARY TABLE temp_review_benefits AS
SELECT id, review_id, benefit FROM review_benefits;

-- Удаляем существующую таблицу бенефитов отзывов
DROP TABLE IF EXISTS review_benefits;

-- Создаем новую таблицу бенефитов отзывов со ссылкой на типы бенефитов
CREATE TABLE IF NOT EXISTS review_benefits (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    benefit_type_id INTEGER NOT NULL REFERENCES benefit_types(id) ON DELETE CASCADE
);

-- Создаем индекс для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_review_benefits_review_id ON review_benefits(review_id);
CREATE INDEX IF NOT EXISTS idx_review_benefits_benefit_type_id ON review_benefits(benefit_type_id);

-- Создаем уникальный индекс, чтобы предотвратить дублирование бенефитов для одного отзыва
CREATE UNIQUE INDEX IF NOT EXISTS idx_review_benefits_unique ON review_benefits(review_id, benefit_type_id);

-- Переносим данные из временной таблицы в новую
-- Попытаемся сопоставить существующие строковые значения с новыми типами бенефитов
WITH benefit_mappings AS (
    SELECT 
        temp.id,
        temp.review_id,
        temp.benefit,
        COALESCE(
            (SELECT bt.id FROM benefit_types bt 
             WHERE temp.benefit ILIKE '%' || bt.name || '%' OR bt.name ILIKE '%' || temp.benefit || '%'
             LIMIT 1),
            1
        ) AS benefit_type_id
    FROM temp_review_benefits temp
)
INSERT INTO review_benefits (review_id, benefit_type_id)
SELECT DISTINCT review_id, benefit_type_id
FROM benefit_mappings;

-- Удаляем временную таблицу
DROP TABLE temp_review_benefits;

-- Записываем информацию о миграции
INSERT INTO migration_history (filename) VALUES ('0009_benefit_types_table.sql'); 