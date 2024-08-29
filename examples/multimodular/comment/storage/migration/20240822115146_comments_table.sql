-- migrate:up
CREATE SCHEMA comment;

CREATE TABLE comment.comment
(
    id   uuid PRIMARY KEY,
    comment text NOT NULL,
    author_id  uuid NOT NULL,
    post_id uuid NOT NULL,
    created_at timestamp  with time zone DEFAULT now() NOT NULL
);


-- migrate:down
DROP TABLE comment.comment;
DROP SCHEMA comment;
