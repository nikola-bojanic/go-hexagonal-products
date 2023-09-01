CREATE TABLE IF NOT EXISTS hex_fwk.product
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    
    name VARCHAR(75) NOT NULL,
    short_description TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(19, 2) DEFAULT 0,
    quantity INT DEFAULT 0,
    category_id BIGINT REFERENCES hex_fwk.category(category_id),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);