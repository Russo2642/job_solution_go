ALTER TABLE companies ADD COLUMN IF NOT EXISTS slug VARCHAR(255);

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

COMMENT ON COLUMN companies.slug IS 'SEO-дружественный URL компании';


UPDATE companies
SET slug = 
    LOWER(
        REGEXP_REPLACE(
            REGEXP_REPLACE(
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

ALTER TABLE companies 
ALTER COLUMN slug SET NOT NULL; 