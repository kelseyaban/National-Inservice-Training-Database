-- Filename: migrations/000022_add_column_to_course_posting.down.sql
ALTER TABLE IF EXISTS course_posting
    DROP CONSTRAINT IF EXISTS fk_course_posting_rank;

ALTER TABLE IF EXISTS course_posting
    DROP COLUMN IF EXISTS rank_id;