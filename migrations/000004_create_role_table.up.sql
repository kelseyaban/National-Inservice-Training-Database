CREATE TABLE IF NOT EXISTS role (
  id bigserial PRIMARY KEY,
  role text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
