-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_items (
    user_id BIGINT NOT NULL,
    sku BIGINT NOT NULL,
    count BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
