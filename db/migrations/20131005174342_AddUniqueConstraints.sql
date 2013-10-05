-- +goose Up
ALTER TABLE movies ADD CONSTRAINT uq_week_title UNIQUE (week, title);
ALTER TABLE services ADD CONSTRAINT uq_movie_id_name UNIQUE (movie_id, name);

-- +goose Down
ALTER TABLE tablename DROP CONSTRAINT uq_week_title;
ALTER TABLE services DROP CONSTRAINT uq_movie_id_name;
