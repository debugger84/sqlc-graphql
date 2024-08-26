-- migrate:up
CREATE SCHEMA post;

CREATE TYPE post.post_status AS ENUM ('draft', 'published', 'deleted');

CREATE TABLE post.post
(
    id   BIGSERIAL PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    status post_status NOT NULL DEFAULT 'draft',
    author_id  uuid NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    published_at TIMESTAMPTZ DEFAULT NULL
);


-- migrate:down
DROP TABLE post.post;
DROP TYPE post.post_status;
DROP SCHEMA post;
