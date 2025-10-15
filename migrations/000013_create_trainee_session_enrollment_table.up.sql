CREATE TABLE IF NOT EXISTS trainee_session_enrollment (
  id bigserial PRIMARY KEY,
  trainee_id bigint NOT NULL REFERENCES trainee(id) ON DELETE CASCADE,
  session_id bigint NOT NULL REFERENCES session(id) ON DELETE CASCADE,
  credithours_completed int NOT NULL DEFAULT 0,
  grade text,
  feedback text,
  version integer NOT NULL DEFAULT 1,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
