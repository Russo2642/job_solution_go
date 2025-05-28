CREATE TABLE IF NOT EXISTS useful_marks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, review_id)
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE indexname = 'idx_useful_marks_user_id' AND tablename = 'useful_marks'
    ) THEN
        CREATE INDEX idx_useful_marks_user_id ON useful_marks(user_id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE indexname = 'idx_useful_marks_review_id' AND tablename = 'useful_marks'
    ) THEN
        CREATE INDEX idx_useful_marks_review_id ON useful_marks(review_id);
    END IF;
END $$;

CREATE OR REPLACE FUNCTION update_review_useful_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE reviews 
        SET useful_count = useful_count + 1
        WHERE id = NEW.review_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE reviews 
        SET useful_count = useful_count - 1
        WHERE id = OLD.review_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DO $$ 
BEGIN
    DROP TRIGGER IF EXISTS trigger_insert_useful_mark ON useful_marks;
    CREATE TRIGGER trigger_insert_useful_mark
    AFTER INSERT ON useful_marks
    FOR EACH ROW
    EXECUTE FUNCTION update_review_useful_count();
    
    DROP TRIGGER IF EXISTS trigger_delete_useful_mark ON useful_marks;
    CREATE TRIGGER trigger_delete_useful_mark
    AFTER DELETE ON useful_marks
    FOR EACH ROW
    EXECUTE FUNCTION update_review_useful_count();
END $$; 