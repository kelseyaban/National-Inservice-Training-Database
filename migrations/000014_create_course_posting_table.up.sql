CREATE TABLE IF NOT EXISTS course_posting (
  id bigserial PRIMARY KEY,
  course_id bigint NOT NULL REFERENCES course(id) ON DELETE CASCADE,
  posting_id bigint NOT NULL REFERENCES posting(id) ON DELETE CASCADE,
  mandatory bool NOT NULL,
  credithours int NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
