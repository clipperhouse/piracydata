-- +goose Up
CREATE TABLE weeks (
    id          serial primary key,
    date		date,
    is_approved boolean
);

-- +goose Down
DROP TABLE weeks;
