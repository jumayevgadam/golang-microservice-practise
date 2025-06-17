-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stock_items (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    sku_id BIGINT NOT NULL UNIQUE REFERENCES sku (sku_id),
    count BIGINT,
    price BIGINT,
    location TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stock_items;
-- +goose StatementEnd
