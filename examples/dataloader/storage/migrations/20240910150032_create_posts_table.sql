-- migrate:up

CREATE TABLE posts
(
    id   BIGSERIAL PRIMARY KEY,
    title text NOT NULL,
    content  text NOT NULL DEFAULT '',
    author_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);


-- migrate:down
DROP TABLE posts;
