CREATE TABLE IF NOT EXISTS course_rank (
  id bigserial PRIMARY KEY,
  course_id bigint NOT NULL REFERENCES course(id) ON DELETE CASCADE,
  rank_id bigint NOT NULL REFERENCES rank(id) ON DELETE CASCADE,
  category_id bigint NOT NULL REFERENCES category(id) ON DELETE CASCADE,
  credithours int NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
