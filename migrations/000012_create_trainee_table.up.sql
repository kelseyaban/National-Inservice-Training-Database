CREATE TABLE IF NOT EXISTS trainee (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  regulation_number text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
