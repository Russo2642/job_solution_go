SET client_min_messages TO WARNING;

CREATE TABLE IF NOT EXISTS rating_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

INSERT INTO rating_categories (name, description) VALUES 
    ('Руководство', 'Качество руководства и менеджмента компании'),
    ('Условия труда', 'Физические условия на рабочем месте, офис, оборудование'),
    ('Коллектив', 'Атмосфера в коллективе, отношения между сотрудниками'),
    ('Уровень дохода', 'Конкурентоспособность заработной платы и финансовых вознаграждений'),
    ('Возможности роста', 'Возможности для карьерного и профессионального развития'),
    ('Условия для отдыха', 'Возможности для отдыха, баланс работы и личной жизни')
ON CONFLICT (name) DO NOTHING;

CREATE TEMPORARY TABLE temp_review_ratings AS
SELECT review_id, category, rating FROM review_category_ratings;

DROP TABLE IF EXISTS review_category_ratings;

CREATE TABLE IF NOT EXISTS review_category_ratings (
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES rating_categories(id) ON DELETE CASCADE,
    rating DECIMAL(3,2) NOT NULL,
    PRIMARY KEY (review_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_review_category_ratings_category_id ON review_category_ratings(category_id);

INSERT INTO review_category_ratings (review_id, category_id, rating)
SELECT temp.review_id, rc.id, temp.rating
FROM temp_review_ratings temp
JOIN rating_categories rc ON rc.name = temp.category;

DROP TABLE temp_review_ratings;

CREATE TEMPORARY TABLE temp_company_ratings AS
SELECT company_id, category, rating FROM company_category_ratings;

DROP TABLE IF EXISTS company_category_ratings;

CREATE TABLE IF NOT EXISTS company_category_ratings (
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES rating_categories(id) ON DELETE CASCADE,
    rating DECIMAL(3,2) NOT NULL,
    PRIMARY KEY (company_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_company_category_ratings_category_id ON company_category_ratings(category_id);

INSERT INTO company_category_ratings (company_id, category_id, rating)
SELECT temp.company_id, rc.id, temp.rating
FROM temp_company_ratings temp
JOIN rating_categories rc ON rc.name = temp.category;

DROP TABLE temp_company_ratings;