-- migrations/001_create_auth_tables.sql

-- 1. Сессии
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    user_agent TEXT,
    ip_address INET,
    expires_at TIMESTAMP NOT NULL,
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- 2. Черный список
CREATE TABLE token_blacklist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_reason VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_blacklist_token ON token_blacklist(token_hash);
CREATE INDEX idx_blacklist_expires ON token_blacklist(expires_at);

-- 3. 2FA
CREATE TABLE two_factor_auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    secret VARCHAR(255) NOT NULL,
    backup_codes TEXT[],
    enabled BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_2fa_user_id ON two_factor_auth(user_id);

-- 4. Временные токены
CREATE TABLE temporary_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_temp_tokens_token ON temporary_tokens(token_hash);
CREATE INDEX idx_temp_tokens_user_type ON temporary_tokens(user_id, type);

-- 5. Аудит
CREATE TABLE auth_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    event_type VARCHAR(50) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN,
    failure_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_user_id ON auth_audit_logs(user_id);
CREATE INDEX idx_audit_created_at ON auth_audit_logs(created_at);

-- Автоматическая очистка истекших записей
CREATE OR REPLACE FUNCTION cleanup_expired_auth_data()
RETURNS void AS $$
BEGIN
    -- Удаляем истекшие сессии
    DELETE FROM sessions WHERE expires_at < NOW();
    
    -- Удаляем истекшие токены из черного списка
    DELETE FROM token_blacklist WHERE expires_at < NOW();
    
    -- Удаляем истекшие временные токены
    DELETE FROM temporary_tokens WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Запуск очистки каждый час (через pg_cron или отдельный worker)