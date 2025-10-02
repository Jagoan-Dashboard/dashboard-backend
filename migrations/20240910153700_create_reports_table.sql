
-- migrations/002_create_reports_table.sql
-- +goose Up
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_name VARCHAR(255) NOT NULL,
    reporter_role VARCHAR(50) NOT NULL,
    village VARCHAR(255) NOT NULL,
    district VARCHAR(255) NOT NULL,
    building_name VARCHAR(255) NOT NULL,
    building_type VARCHAR(50) NOT NULL,
    report_status VARCHAR(50) NOT NULL,
    funding_source VARCHAR(50) NOT NULL,
    last_year_construction INTEGER,
    full_address TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    floor_area DECIMAL(10, 2),
    floor_count INTEGER,
    work_type VARCHAR(50),
    condition_after_rehab VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reports_village ON reports(village);
CREATE INDEX idx_reports_district ON reports(district);
CREATE INDEX idx_reports_building_type ON reports(building_type);

-- +goose Down
DROP TABLE IF EXISTS reports;