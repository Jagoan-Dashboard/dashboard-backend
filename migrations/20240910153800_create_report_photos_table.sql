-- migrations/003_create_report_photos_table.sql
-- +goose Up
CREATE TABLE report_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_report_photos_report_id ON report_photos(report_id);

-- +goose Down
DROP TABLE IF EXISTS report_photos;