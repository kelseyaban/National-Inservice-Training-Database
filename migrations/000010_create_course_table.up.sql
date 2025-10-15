CREATE TABLE IF NOT EXISTS course (
  id bigserial PRIMARY KEY,
  course text NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
