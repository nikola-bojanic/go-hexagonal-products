CREATE TABLE IF NOT EXISTS hex_fwk.order
(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    
    status VARCHAR(75),
    user_id UUID REFERENCES hex_fwk.user(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);