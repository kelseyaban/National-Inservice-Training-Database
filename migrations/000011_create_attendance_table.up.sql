CREATE TABLE IF NOT EXISTS attendance (
  id bigserial PRIMARY KEY,
  user_session_id bigint NOT NULL REFERENCES user_session(id) ON  DELETE CASCADE,
  attendance bool NOT NULL,
  date date NOT NULL,
 created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()

);
