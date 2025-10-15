CREATE TABLE IF NOT EXISTS posting (
  id bigserial PRIMARY KEY,
  posting text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
