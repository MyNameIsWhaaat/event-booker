CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT        NOT NULL,
    starts_at           TIMESTAMPTZ NOT NULL,
    capacity            INT         NOT NULL CHECK (capacity > 0),
    requires_payment    BOOLEAN     NOT NULL DEFAULT TRUE,
    booking_ttl_seconds INT         NOT NULL CHECK (booking_ttl_seconds > 0),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS events_starts_at_idx
ON events (starts_at);

CREATE TABLE IF NOT EXISTS bookings (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id     UUID        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_email   TEXT        NOT NULL,
    status       TEXT        NOT NULL CHECK (status IN ('pending', 'confirmed', 'cancelled')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL,
    confirmed_at TIMESTAMPTZ NULL,
    cancelled_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS bookings_event_status_idx
ON bookings (event_id, status);

CREATE INDEX IF NOT EXISTS bookings_status_expires_idx
ON bookings (status, expires_at);