DROP INDEX IF EXISTS bookings_status_expires_idx;
DROP INDEX IF EXISTS bookings_event_status_idx;
DROP TABLE IF EXISTS bookings;

DROP INDEX IF EXISTS events_starts_at_idx;
DROP TABLE IF EXISTS events;

DROP EXTENSION IF EXISTS pgcrypto;