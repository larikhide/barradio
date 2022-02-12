

-- DROP TABLE IF EXISTS vote;

CREATE TABLE IF NOT EXISTS vote
(
    id uuid NOT NULL,
    category_id uuid NOT NULL,
    category_code integer NOT NULL,
    date timestamp NOT NULL,
    CONSTRAINT vote_pkey PRIMARY KEY (id)
);



COMMENT ON TABLE vote IS 'голоса';
COMMENT ON COLUMN vote.id  IS 'идентификатор';
COMMENT ON COLUMN vote.category_id IS 'иднгтификатор категории (uuid)';
COMMENT ON COLUMN vote.category_code IS 'код категории (int)';
COMMENT ON COLUMN vote.date  IS 'дата время голосования';


