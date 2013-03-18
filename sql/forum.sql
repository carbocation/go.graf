-- SQL file for the app

CREATE TABLE user (
	id BIGSERIAL PRIMARY KEY,
	handle TEXT
)

CREATE TABLE entry (
    id BIGSERIAL PRIMARY KEY,
    title TEXT,
    body TEXT,
	created TIMESTAMP WITH TIME ZONE,
    author_id BIGINT REFERENCES user(id)
);
