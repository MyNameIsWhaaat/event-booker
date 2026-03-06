DROP INDEX IF EXISTS bookings_event_user_active_uidx;
DROP INDEX IF EXISTS bookings_user_id_idx;

ALTER TABLE bookings
DROP CONSTRAINT IF EXISTS bookings_user_id_fkey;

ALTER TABLE bookings
DROP COLUMN IF EXISTS user_id;

DROP TABLE IF EXISTS users;