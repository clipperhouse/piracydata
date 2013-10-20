-- +goose Up
INSERT INTO weeks (date, is_approved)
	SELECT DISTINCT week, TRUE
	FROM movies;

-- +goose Down
TRUNCATE TABLE weeks;
