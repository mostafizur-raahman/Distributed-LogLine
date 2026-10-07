CREATE TABLE IF NOT EXISTS logs (
    id         BIGSERIAL    PRIMARY KEY,
    level      TEXT         NOT NULL CHECK (level IN ('debug', 'info', 'warn', 'error', 'fatal')),
    message    TEXT         NOT NULL CHECK (length(message) <= 4096),
    service    TEXT         NOT NULL CHECK (length(service) <= 128),
    timestamp  TIMESTAMPTZ  NOT NULL,
    data       JSONB        DEFAULT '{}',
    created_at TIMESTAMPTZ  DEFAULT NOW()
);