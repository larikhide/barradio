

-- DROP TABLE IF EXISTS vote;

CREATE TABLE IF NOT EXISTS vote
(
    id uuid NOT NULL,
    category_id uuid NOT NULL,
    category_code integer NOT NULL,
    created_at timestamp NOT NULL,
    CONSTRAINT vote_pkey PRIMARY KEY (id)
);



COMMENT ON TABLE vote IS 'голоса';
COMMENT ON COLUMN vote.id  IS 'идентификатор';
COMMENT ON COLUMN vote.category_id IS 'идентификатор категории (uuid)';
COMMENT ON COLUMN vote.category_code IS 'код категории (int)';
COMMENT ON COLUMN vote.created_at  IS 'дата время голосования';


