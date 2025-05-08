-- Добавление поля is_former_employee в таблицу reviews
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS is_former_employee BOOLEAN NOT NULL DEFAULT FALSE;

-- Добавление индекса для нового поля
CREATE INDEX IF NOT EXISTS idx_reviews_is_former_employee ON reviews(is_former_employee);