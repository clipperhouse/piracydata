-- +goose Up
CREATE TABLE services (
    id          serial primary key,
    movie_id    integer,
    name	    varchar(50),
    available   boolean
);

-- +goose Down
DROP TABLE services;
