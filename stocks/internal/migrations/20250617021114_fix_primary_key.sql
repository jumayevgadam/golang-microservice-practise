-- +goose Up
-- +goose StatementBegin
ALTER TABLE stock_items
ADD COLUMN stock_item_id BIGSERIAL PRIMARY KEY;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stock_items;
DROP TABLE IF EXISTS sku;
-- +goose StatementEnd
