-- Database schema for brute-force protected login system

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Failed login attempts tracking (per user)
CREATE TABLE IF NOT EXISTS user_failed_attempts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45) NOT NULL
);

-- IP-based failed attempts tracking
CREATE TABLE IF NOT EXISTS ip_failed_attempts (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(45) NOT NULL,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email VARCHAR(255) NOT NULL
);

-- User suspensions
CREATE TABLE IF NOT EXISTS user_suspensions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    suspended_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- IP blocks
CREATE TABLE IF NOT EXISTS ip_blocks (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(45) UNIQUE NOT NULL,
    blocked_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_failed_attempts_email_time ON user_failed_attempts(email, attempted_at);
CREATE INDEX IF NOT EXISTS idx_user_failed_attempts_time ON user_failed_attempts(attempted_at);
CREATE INDEX IF NOT EXISTS idx_ip_failed_attempts_ip_time ON ip_failed_attempts(ip_address, attempted_at);
CREATE INDEX IF NOT EXISTS idx_ip_failed_attempts_time ON ip_failed_attempts(attempted_at);
CREATE INDEX IF NOT EXISTS idx_user_suspensions_email ON user_suspensions(email);
CREATE INDEX IF NOT EXISTS idx_ip_blocks_ip ON ip_blocks(ip_address);

-- Cleanup function to remove old records (optional, can be run periodically)
-- DELETE FROM user_failed_attempts WHERE attempted_at < NOW() - INTERVAL '1 hour';
-- DELETE FROM ip_failed_attempts WHERE attempted_at < NOW() - INTERVAL '1 hour';

