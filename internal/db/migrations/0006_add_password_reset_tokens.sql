-- Создаем таблицу для токенов сброса пароля
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Добавляем индекс для токена и индекс для user_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_indexes 
        WHERE indexname = 'idx_password_reset_tokens_token' AND tablename = 'password_reset_tokens'
    ) THEN
        CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_indexes 
        WHERE indexname = 'idx_password_reset_tokens_user_id' AND tablename = 'password_reset_tokens'
    ) THEN
        CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
    END IF;
END $$;

-- Добавляем комментарии к таблице и полям
COMMENT ON TABLE password_reset_tokens IS 'Токены для сброса пароля пользователей';
COMMENT ON COLUMN password_reset_tokens.id IS 'Уникальный идентификатор токена';
COMMENT ON COLUMN password_reset_tokens.user_id IS 'ID пользователя, которому принадлежит токен';
COMMENT ON COLUMN password_reset_tokens.token IS 'Значение токена';
COMMENT ON COLUMN password_reset_tokens.expires_at IS 'Время истечения срока действия токена';
COMMENT ON COLUMN password_reset_tokens.created_at IS 'Время создания токена'; 