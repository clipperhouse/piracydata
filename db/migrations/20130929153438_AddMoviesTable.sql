-- +goose Up
CREATE TABLE movies (
    id          serial primary key,
    week        date NOT NULL,
    title       varchar(200),
    imdb        varchar(20),
    rank        integer,
    stream      boolean,
    rent        boolean,
    buy         boolean
);

-- +goose Down
DROP TABLE movies;
