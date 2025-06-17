-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items 
ADD CONSTRAINT cart_items_pkey PRIMARY KEY (user_id, sku);;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cart_items 
DROP CONSTRAINT cart_items_pkey;
-- +goose StatementEnd
