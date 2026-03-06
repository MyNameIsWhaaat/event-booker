CREATE TABLE IF NOT EXISTS users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE bookings
ADD COLUMN IF NOT EXISTS user_id UUID;

INSERT INTO users (email)
SELECT DISTINCT b.user_email
FROM bookings b
WHERE b.user_email IS NOT NULL
  AND b.user_email <> ''
ON CONFLICT (email) DO NOTHING;

UPDATE bookings b
SET user_id = u.id
FROM users u
WHERE b.user_email = u.email
  AND b.user_id IS NULL;

ALTER TABLE bookings
ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE bookings
ADD CONSTRAINT bookings_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS bookings_user_id_idx
ON bookings (user_id);

CREATE UNIQUE INDEX IF NOT EXISTS bookings_event_user_active_uidx
ON bookings (event_id, user_id)
WHERE status IN ('pending', 'confirmed');