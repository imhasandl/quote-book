-- +goose Up
CREATE TABLE quotes(
   id SERIAL PRIMARY KEY,
   author TEXT NOT NULL,
   quote TEXT NOT NULL
);

-- +goose Down
DROP TABLE quotes;