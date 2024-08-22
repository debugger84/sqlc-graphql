-- migrate:up
CREATE TYPE author_status AS ENUM ('active', 'inactive', 'deleted');

CREATE TABLE authors
(
    id   BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    bio  text,
    status author_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT now()
);


-- migrate:down
DROP TABLE authors;
DROP TYPE author_status;
