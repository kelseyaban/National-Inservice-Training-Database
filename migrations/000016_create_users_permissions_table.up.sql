CREATE TABLE IF NOT EXISTS users_permissions (
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  permission_id bigint NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, permission_id)
);
