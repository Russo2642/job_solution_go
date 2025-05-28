ALTER TABLE reviews ADD COLUMN IF NOT EXISTS is_recommended BOOLEAN DEFAULT TRUE;

ALTER TABLE companies ADD COLUMN IF NOT EXISTS recommendation_percentage DECIMAL(5,2) DEFAULT 0;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_description 
        WHERE objoid = 'reviews'::regclass::oid 
        AND objsubid = (
            SELECT attnum FROM pg_attribute 
            WHERE attrelid = 'reviews'::regclass 
            AND attname = 'is_recommended'
        )
    ) THEN
        COMMENT ON COLUMN reviews.is_recommended IS 'Флаг, указывающий рекомендует ли пользователь компанию';
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_description 
        WHERE objoid = 'companies'::regclass::oid 
        AND objsubid = (
            SELECT attnum FROM pg_attribute 
            WHERE attrelid = 'companies'::regclass 
            AND attname = 'recommendation_percentage'
        )
    ) THEN
        COMMENT ON COLUMN companies.recommendation_percentage IS 'Процент пользователей, рекомендующих компанию';
    END IF;
END $$;

UPDATE companies c
SET recommendation_percentage = 100
WHERE EXISTS (
    SELECT 1 FROM reviews r 
    WHERE r.company_id = c.id AND r.status = 'approved'
); 