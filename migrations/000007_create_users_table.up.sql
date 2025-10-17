CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  regulation_number text NOT NULL,
  fname text NOT NULL,
  lname text NOT NULL,
  email citext NOT NULL UNIQUE,
  password_hash bytea NOT NULL,
  activated bool NOT NULL DEFAULT false,
  gender_id char(1) NOT NULL,
  formation_id bigint NOT NULL REFERENCES formation(id) ON DELETE RESTRICT,
  rank_id bigint REFERENCES rank(id) ON DELETE SET NULL,
  posting_id bigint REFERENCES posting(id) ON DELETE SET NULL,
  version integer NOT NULL DEFAULT 1,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
