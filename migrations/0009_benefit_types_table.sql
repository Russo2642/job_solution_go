SET client_min_messages TO WARNING;

CREATE TABLE IF NOT EXISTS benefit_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

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

CREATE TEMPORARY TABLE temp_review_benefits AS
SELECT id, review_id, benefit FROM review_benefits;

DROP TABLE IF EXISTS review_benefits;

CREATE TABLE IF NOT EXISTS review_benefits (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    benefit_type_id INTEGER NOT NULL REFERENCES benefit_types(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_review_benefits_review_id ON review_benefits(review_id);
CREATE INDEX IF NOT EXISTS idx_review_benefits_benefit_type_id ON review_benefits(benefit_type_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_review_benefits_unique ON review_benefits(review_id, benefit_type_id);

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

DROP TABLE temp_review_benefits;