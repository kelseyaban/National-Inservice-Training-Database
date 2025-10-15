CREATE TABLE IF NOT EXISTS attendance (
  id bigserial PRIMARY KEY,
  trainee_course_id bigint NOT NULL REFERENCES trainee_session_enrollment(id) ON DELETE CASCADE,
  attendance bool NOT NULL,
  date date NOT NULL,
 created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()

);
