CREATE TABLE IF NOT EXISTS formation (
  id bigserial PRIMARY KEY,
  region_id bigint NOT NULL REFERENCES region(id) ON DELETE CASCADE,
  formation text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
