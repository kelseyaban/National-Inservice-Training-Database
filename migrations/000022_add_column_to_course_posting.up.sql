ALTER TABLE course_posting
    ADD COLUMN rank_id bigint NOT NULL;

ALTER TABLE course_posting
    ADD CONSTRAINT fk_course_posting_rank
        FOREIGN KEY (rank_id) REFERENCES rank(id)
        ON DELETE CASCADE;
