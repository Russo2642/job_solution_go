-- Добавляем поле city_id в таблицу companies
ALTER TABLE companies ADD COLUMN IF NOT EXISTS city_id INTEGER;

-- Добавляем внешний ключ на таблицу cities, если его еще нет
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_companies_city' AND conrelid = 'companies'::regclass
    ) THEN
        ALTER TABLE companies ADD CONSTRAINT fk_companies_city 
            FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Создаем индекс для ускорения поиска по city_id, если его еще нет
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE indexname = 'idx_companies_city_id' AND tablename = 'companies'
    ) THEN
        CREATE INDEX idx_companies_city_id ON companies(city_id);
    END IF;
END $$;

-- Добавим комментарий к таблице, если его еще нет
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_description 
        WHERE objoid = 'companies'::regclass::oid 
        AND objsubid = (
            SELECT attnum FROM pg_attribute 
            WHERE attrelid = 'companies'::regclass 
            AND attname = 'city_id'
        )
    ) THEN
        COMMENT ON COLUMN companies.city_id IS 'ID города из справочника';
    END IF;
END $$; 