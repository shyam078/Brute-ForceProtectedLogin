-- Initialize database with sample users
-- Password for alice@example.com: password123
-- Password for bob@example.com: password123
-- Password hash is bcrypt hash of "password123"

INSERT INTO users (email, password_hash) VALUES
('alice@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
('bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy')
ON CONFLICT (email) DO NOTHING;

