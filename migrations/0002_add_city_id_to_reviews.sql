ALTER TABLE reviews ADD COLUMN IF NOT EXISTS city_id INTEGER;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_reviews_city' AND conrelid = 'reviews'::regclass
    ) THEN
        ALTER TABLE reviews ADD CONSTRAINT fk_reviews_city 
            FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE indexname = 'idx_reviews_city_id' AND tablename = 'reviews'
    ) THEN
        CREATE INDEX idx_reviews_city_id ON reviews(city_id);
    END IF;
END $$;

UPDATE reviews r
SET city_id = (
    SELECT c.id 
    FROM cities c 
    WHERE c.name = r.city
    LIMIT 1
)
WHERE r.city_id IS NULL AND r.city IS NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_description 
        WHERE objoid = 'reviews'::regclass::oid 
        AND objsubid = (
            SELECT attnum FROM pg_attribute 
            WHERE attrelid = 'reviews'::regclass 
            AND attname = 'city_id'
        )
    ) THEN
        COMMENT ON COLUMN reviews.city_id IS 'ID города из справочника';
    END IF;
END $$; 