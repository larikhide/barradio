-- DROP TABLE IF EXISTS composition;

CREATE TABLE IF NOT EXISTS composition
(
    id uuid NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    categoru_code integer NOT NULL,
    category_id uuid NOT NULL,
    singer character varying(256) NOT NULL DEFAULT '',
    image_url text NOT NULL DEFAULT 'no image',
    composition_url text NOT NULL DEFAULT 'no song',
    duration timestamp NOT NULL,
    is_run boolean NOT NULL DEFAULT true,
    CONSTRAINT composition_pkey PRIMARY KEY (id)
);

COMMENT ON COLUMN composition.id  IS 'Идентификатор';
COMMENT ON COLUMN composition.categoru_code IS 'код категории (int)';
COMMENT ON COLUMN composition.category_id IS 'идентификатор категории (uuid)';
COMMENT ON COLUMN composition.singer IS 'исполнитель';
COMMENT ON COLUMN composition.image_url IS 'путь к картинке';
COMMENT ON COLUMN composition.composition_url IS 'путь к песне';
COMMENT ON COLUMN composition.created_at IS 'дата и время создания';
COMMENT ON COLUMN composition.updated_at IS 'дата и время обноления';
COMMENT ON COLUMN composition.is_run  IS 'признак рабочей записи';
COMMENT ON COLUMN composition.duration IS 'время звучания';




