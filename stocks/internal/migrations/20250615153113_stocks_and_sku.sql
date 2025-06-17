-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sku (
    sku_id BIGINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT
);

INSERT INTO sku (sku_id, name, type) VALUES
(1001, 't-shirt', 'apparel'),
(2020, 'cup', 'accessory'),
(3033, 'book', 'stationery'),
(4044, 'pen', 'stationery'),
(5055, 'powerbank', 'electronics'),
(6066, 'hoody', 'apparel'),
(7077, 'umbrella', 'accessory'),
(8088, 'socks', 'apparel'),
(9099, 'wallet', 'accessory'),
(10101, 'pink-hoody', 'apparel');

CREATE TABLE IF NOT EXISTS stock_items (
    user_id BIGINT NOT NULL,
    sku BIGINT NOT NULL,
    count BIGINT,
    name TEXT,
    type TEXT,
    price BIGINT,
    location TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stock_items;
DROP TABLE IF EXISTS sku;
-- +goose StatementEnd
