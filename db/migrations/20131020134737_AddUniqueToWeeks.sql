-- +goose Up
ALTER TABLE weeks ADD CONSTRAINT uq_week_date UNIQUE (date);

-- +goose Down
ALTER TABLE weeks DROP CONSTRAINT uq_week_date;
