-- Initialize database with sample users
-- Password for alice@example.com: password123
-- Password for bob@example.com: password123
-- Password hash is bcrypt hash of "password123"

INSERT INTO users (email, password_hash) VALUES
('alice@example.com', '$2a$10$ANQwEd9x9uUf1AKfilZ8SOZwfW2aqHrF6jvC9sfcjEBx4Wy6WiwXq'),
('bob@example.com', '$2a$10$ANQwEd9x9uUf1AKfilZ8SOZwfW2aqHrF6jvC9sfcjEBx4Wy6WiwXq')
ON CONFLICT (email) DO NOTHING;

