-- +goose Up
CREATE TABLE movies (
    id          serial primary key,
    week		date,
    title       varchar(200),
    imdb        varchar(50),
    rank        integer,
    streaming   integer,
    rental      integer,
    purchase    integer,
    dvd         integer
);

-- +goose Down
DROP TABLE movies;
