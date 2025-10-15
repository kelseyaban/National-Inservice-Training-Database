CREATE TABLE IF NOT EXISTS users_role (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id bigint NOT NULL REFERENCES role(id) ON DELETE CASCADE,
  UNIQUE(user_id, role_id)
);
