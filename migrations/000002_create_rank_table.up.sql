CREATE TABLE IF NOT EXISTS rank (
  id bigserial PRIMARY KEY,
  title text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
