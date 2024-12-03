CREATE TABLE IF NOT EXISTS urls (
    "id" BIGSERIAL PRIMARY KEY,
    "orignal_url" TEXT NOT NULL,
    "short_code" TEXT NOT NULL UNIQUE,
    "is_custom" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expired_at" TIMESTAMP NOT NULL
);

CREATE INDEX idx_short_code ON urls(short_code);
CREATE INDEX idx_expired_at ON urls(expired_at);