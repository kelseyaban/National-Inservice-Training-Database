CREATE TABLE IF NOT EXISTS session (
  id bigserial PRIMARY KEY,
  course_id bigint NOT NULL REFERENCES course(id) ON DELETE CASCADE,
  formation_id bigint NOT NULL REFERENCES formation(id) ON DELETE CASCADE,
  facilitator_id bigint REFERENCES users(id) ON DELETE SET NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
