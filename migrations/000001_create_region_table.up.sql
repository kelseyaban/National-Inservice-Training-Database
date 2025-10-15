CREATE TABLE IF NOT EXISTS region (
  id bigserial PRIMARY KEY,
  region text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
