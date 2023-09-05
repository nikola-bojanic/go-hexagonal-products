CREATE TABLE IF NOT EXISTS hex_fwk.order_product
(
    order_id UUID DEFAULT uuid_generate_v4(),
    product_id BIGINT,
    PRIMARY KEY(order_id, product_id),
    quantity INT,
    FOREIGN KEY(order_id) REFERENCES hex_fwk.order(id),
    FOREIGN KEY(product_id) REFERENCES hex_fwk.product(id)
);