-- Включаем безопасное удаление (без ошибок при отсутствии объектов)
SET client_min_messages TO WARNING;

-- Создаем таблицу для отслеживания миграций, если её еще нет
CREATE TABLE IF NOT EXISTS migration_history (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создаем типы для перечислений
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('user', 'moderator', 'admin');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'review_status') THEN
        CREATE TYPE review_status AS ENUM ('pending', 'approved', 'rejected');
    END IF;
END $$;

-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создаем таблицу городов
CREATE TABLE IF NOT EXISTS cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    UNIQUE(name, region, country)
);

-- Создаем таблицу компаний
CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    size VARCHAR(50) NOT NULL,
    logo VARCHAR(255),
    website VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    reviews_count INTEGER NOT NULL DEFAULT 0,
    average_rating DECIMAL(3,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создаем таблицу отраслей
CREATE TABLE IF NOT EXISTS industries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Создаем таблицу для связи компаний с отраслями
CREATE TABLE IF NOT EXISTS company_industries (
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    industry_id INTEGER NOT NULL REFERENCES industries(id) ON DELETE CASCADE,
    PRIMARY KEY (company_id, industry_id)
);

-- Создаем таблицу для рейтингов компаний по категориям
CREATE TABLE IF NOT EXISTS company_category_ratings (
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    category VARCHAR(50) NOT NULL,
    rating DECIMAL(3,2) NOT NULL,
    PRIMARY KEY (company_id, category)
);

-- Создаем таблицу отзывов
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    position VARCHAR(100) NOT NULL,
    employment VARCHAR(50) NOT NULL,
    employment_period VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    rating DECIMAL(3,2) NOT NULL,
    pros TEXT NOT NULL,
    cons TEXT NOT NULL,
    status review_status NOT NULL DEFAULT 'pending',
    moderation_comment TEXT,
    useful_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMP
);

-- Создаем таблицу для рейтингов отзывов по категориям
CREATE TABLE IF NOT EXISTS review_category_ratings (
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    category VARCHAR(50) NOT NULL,
    rating DECIMAL(3,2) NOT NULL,
    PRIMARY KEY (review_id, category)
);

-- Создаем таблицу для бонусов и льгот в отзывах
CREATE TABLE IF NOT EXISTS review_benefits (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    benefit VARCHAR(100) NOT NULL
);

-- Создаем таблицу для refresh токенов
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Добавляем основные отрасли IT
INSERT INTO industries (name) VALUES 
    ('Разработка ПО'),
    ('Веб-разработка'),
    ('Мобильная разработка'),
    ('Искусственный интеллект и машинное обучение'),
    ('Большие данные и аналитика'),
    ('Кибербезопасность'),
    ('Облачные вычисления'),
    ('DevOps и SRE'),
    ('Финтех'),
    ('Блокчейн'),
    ('Игровая индустрия'),
    ('E-commerce'),
    ('Интернет вещей (IoT)'),
    ('Телекоммуникации'),
    ('ERP/CRM системы')
ON CONFLICT (name) DO NOTHING;

-- Создаем индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_reviews_company_id ON reviews(company_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_companies_name ON companies(name);
CREATE INDEX IF NOT EXISTS idx_company_industries_company_id ON company_industries(company_id);
CREATE INDEX IF NOT EXISTS idx_company_industries_industry_id ON company_industries(industry_id);
CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);
CREATE INDEX IF NOT EXISTS idx_cities_country ON cities(country);

-- Функция для обновления updated_at при обновлении записи
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры для обновления timestamps
DO $$ 
BEGIN
    -- Триггер для пользователей
    DROP TRIGGER IF EXISTS update_users_updated_at ON users;
    CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
    
    -- Триггер для компаний
    DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;
    CREATE TRIGGER update_companies_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
    
    -- Триггер для отзывов
    DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
    CREATE TRIGGER update_reviews_updated_at
    BEFORE UPDATE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
END $$; 

-- Добавляем базовые города Казахстана
INSERT INTO cities (name, region, country)
VALUES 
    ('Алматы', 'Алматинская область', 'Казахстан'),
    ('Астана', 'Акмолинская область', 'Казахстан'),
    ('Шымкент', 'Туркестанская область', 'Казахстан'),
    ('Караганда', 'Карагандинская область', 'Казахстан'),
    ('Актобе', 'Актюбинская область', 'Казахстан'),
    ('Тараз', 'Жамбылская область', 'Казахстан'),
    ('Павлодар', 'Павлодарская область', 'Казахстан'),
    ('Усть-Каменогорск', 'Восточно-Казахстанская область', 'Казахстан'),
    ('Семей', 'Восточно-Казахстанская область', 'Казахстан'),
    ('Атырау', 'Атырауская область', 'Казахстан'),
    ('Костанай', 'Костанайская область', 'Казахстан'),
    ('Кызылорда', 'Кызылординская область', 'Казахстан'),
    ('Уральск', 'Западно-Казахстанская область', 'Казахстан'),
    ('Петропавловск', 'Северо-Казахстанская область', 'Казахстан'),
    ('Кокшетау', 'Акмолинская область', 'Казахстан'),
    ('Талдыкорган', 'Алматинская область', 'Казахстан'),
    ('Экибастуз', 'Павлодарская область', 'Казахстан'),
    ('Туркестан', 'Туркестанская область', 'Казахстан')
ON CONFLICT (name, region, country) DO NOTHING; 