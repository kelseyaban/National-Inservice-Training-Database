CREATE TABLE IF NOT EXISTS category (
  id bigserial PRIMARY KEY,
  category text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
