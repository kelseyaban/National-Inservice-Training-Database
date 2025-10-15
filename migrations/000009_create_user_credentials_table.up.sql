CREATE TABLE IF NOT EXISTS user_credentials (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  email citext NOT NULL UNIQUE,
  password_hash bytea NOT NULL,
  activated bool NOT NULL DEFAULT false,
  version integer NOT NULL DEFAULT 1
);
