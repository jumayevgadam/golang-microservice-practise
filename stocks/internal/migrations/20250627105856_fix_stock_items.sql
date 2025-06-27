-- +goose Up
-- +goose StatementBegin
ALTER TABLE stock_items DROP CONSTRAINT stock_items_sku_id_key;
ALTER TABLE stock_items ADD CONSTRAINT unique_user_sku UNIQUE (user_id, sku_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stock_items DROP CONSTRAINT unique_user_sku;
ALTER TABLE stock_items ADD CONSTRAINT stock_items_sku_id_key UNIQUE (sku_id);
-- +goose StatementEnd
