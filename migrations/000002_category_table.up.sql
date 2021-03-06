--DROP TABLE category;

CREATE TABLE IF NOT EXISTS category
(
    id uuid NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name character varying(50) NOT NULL,
    code integer generated always as identity,
    name_ru character varying(50)  NOT NULL,
    is_run boolean NOT NULL DEFAULT true,
    CONSTRAINT pk_category_id PRIMARY KEY (id)
);

COMMENT ON COLUMN category.id IS 'Идентификатор';
COMMENT ON COLUMN category.created_at  IS 'Дата и время создания';
COMMENT ON COLUMN category.updated_at IS 'Дата и время изменения';
COMMENT ON COLUMN category.name IS 'Название категории на англиском';
COMMENT ON COLUMN category.code IS 'цифрофой код категории (на всякий случай)';
COMMENT ON COLUMN category.name_ru IS 'Название категории на русском';
COMMENT ON COLUMN category.is_run IS 'Признак рабочей записи';


INSERT into category(id, created_at, updated_at, name, name_ru )
       VALUES(uuid_generate_v4(), now(), now(),'cheerful', 'веселье'),
            (uuid_generate_v4(), now(), now(),'relaxed', 'релах'),
            (uuid_generate_v4(), now(), now(),'lyrical', 'лирика');