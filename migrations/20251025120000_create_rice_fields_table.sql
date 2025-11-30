-- migrations/20251025120000_create_rice_fields_table.sql
-- +goose Up
CREATE TABLE rice_fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    district VARCHAR(255) NOT NULL,
    longitude DECIMAL(12,9),
    latitude DECIMAL(12,9),
    date DATE NOT NULL,
    rainfed_rice_fields DECIMAL(15, 2),
    irrigated_rice_fields DECIMAL(15, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rice_fields_district ON rice_fields(district);
CREATE INDEX idx_rice_fields_date ON rice_fields(date);

-- +goose Down
DROP TABLE IF EXISTS rice_fields;