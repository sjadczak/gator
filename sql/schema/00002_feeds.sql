-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR NOT NULL,
    url VARCHAR UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd

-- Like any table in our DB, we'll need the standard id, created_at, and updated_at fields. We'll also need a few more:

-- name: The name of the feed (like "The Changelog, or "The Boot.dev Blog")
-- url: The URL of the feed
-- user_id: The ID of the user who added this feed
-- Make the url field unique so that in the future we aren't downloading duplicate posts.

-- Use an ON DELETE CASCADE constraint on the user_id foreign key so that if a user is deleted, all of their 
-- feeds are automatically deleted as well. This will ensure we have no orphaned records and that deleting 
-- the users in the reset command also deletes all of their feeds.
