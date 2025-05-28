SET client_min_messages TO WARNING;

CREATE TYPE suggestion_type AS ENUM ('company', 'suggestion');

CREATE TABLE IF NOT EXISTS suggestions (
    id SERIAL PRIMARY KEY,
    type suggestion_type NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_suggestions_type ON suggestions(type);
CREATE INDEX IF NOT EXISTS idx_suggestions_created_at ON suggestions(created_at); 