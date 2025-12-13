-- ============================================================
-- AUTHENTICATION SYSTEM SCHEMA (POSTGRESQL - 2025 BEST PRACTICE)
-- Using BIGINT auto-increment IDs (BIGSERIAL)
-- ============================================================


-- ======================================
-- TABLE: users
-- ======================================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    password_hash TEXT NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    -- CONSTRAINT role_check CHECK (role IN ('user', 'admin'))
);

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);



-- ======================================
-- TABLE: refresh_token_families
-- ======================================
CREATE TABLE refresh_token_families (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_rtf_user_id ON refresh_token_families (user_id);



-- ======================================
-- TABLE: refresh_tokens
-- ======================================
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_id BIGINT NOT NULL REFERENCES refresh_token_families(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    ip_address VARCHAR(60),
    user_agent TEXT
);

CREATE INDEX idx_rt_user_family ON refresh_tokens (user_id, family_id);
CREATE INDEX idx_rt_expires_at ON refresh_tokens (expires_at);



-- ======================================
-- TABLE: password_reset_tokens
-- ======================================
CREATE TABLE password_reset_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_prt_user_id ON password_reset_tokens (user_id);
CREATE INDEX idx_prt_expires_at ON password_reset_tokens (expires_at);



-- ======================================
-- TABLE: email_verification_tokens
-- ======================================
CREATE TABLE email_verification_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_evt_user_id ON email_verification_tokens (user_id);
CREATE INDEX idx_evt_expires_at ON email_verification_tokens (expires_at);



-- ======================================
-- TABLE: login_attempts (optional, security)
-- ======================================
CREATE TABLE login_attempts (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255),
    ip_address VARCHAR(60),
    user_agent TEXT,
    success BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_la_email ON login_attempts (email);
CREATE INDEX idx_la_ip ON login_attempts (ip_address);
CREATE INDEX idx_la_created_at ON login_attempts (created_at);



-- ======================================
-- TABLE: user_sessions (optional)
-- ======================================
CREATE TABLE user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address VARCHAR(60),
    user_agent TEXT,
    login_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ,
    logout_at TIMESTAMPTZ
);

CREATE INDEX idx_us_user_id ON user_sessions (user_id);
CREATE INDEX idx_us_last_seen ON user_sessions (last_seen_at);

