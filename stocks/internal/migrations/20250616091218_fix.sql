-- +goose Up
-- +goose StatementBegin
ALTER TABLE stock_items
ADD COLUMN sku_id BIGINT REFERENCES sku(sku_id);

ALTER TABLE stock_items
DROP COLUMN sku;

ALTER TABLE stock_items
DROP COLUMN name;

ALTER TABLE stock_items
DROP COLUMN type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stock_items;
DROP TABLE IF EXISTS sku;
-- +goose StatementEnd
