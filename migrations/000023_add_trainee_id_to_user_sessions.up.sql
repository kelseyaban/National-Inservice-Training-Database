ALTER TABLE user_session
ADD COLUMN trainee_id BIGINT;

ALTER TABLE user_session
ADD CONSTRAINT fk_user_session_user
FOREIGN KEY (trainee_id) REFERENCES users(id) ON DELETE CASCADE;