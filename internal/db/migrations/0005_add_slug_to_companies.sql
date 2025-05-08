-- Добавляем поле slug в таблицу companies
ALTER TABLE companies ADD COLUMN IF NOT EXISTS slug VARCHAR(255);

-- Создаем уникальный индекс для slug
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_indexes 
        WHERE indexname = 'idx_companies_slug' AND tablename = 'companies'
    ) THEN
        CREATE UNIQUE INDEX idx_companies_slug ON companies(slug);
    END IF;
END $$;

-- Добавляем комментарий к полю slug
COMMENT ON COLUMN companies.slug IS 'SEO-дружественный URL компании';

-- Заполняем slug для существующих записей
-- Мы используем комбинацию name и id для обеспечения уникальности
UPDATE companies
SET slug = 
    LOWER(
        REGEXP_REPLACE(
            REGEXP_REPLACE(
                -- Транслитерация не поддерживается в стандартном PostgreSQL, 
                -- поэтому это будет делаться в Go-коде
                name,
                '[^a-zA-Z0-9\s]', 
                '-', 
                'g'
            ),
            '[-\s]+', 
            '-', 
            'g'
        )
    ) || '-' || id
WHERE slug IS NULL;

-- Устанавливаем ограничение NOT NULL для slug после заполнения данных
ALTER TABLE companies 
ALTER COLUMN slug SET NOT NULL; 