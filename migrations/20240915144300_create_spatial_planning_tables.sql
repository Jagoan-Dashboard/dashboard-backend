-- migrations/004_create_spatial_planning_tables.sql
-- +goose Up
CREATE TABLE spatial_planning_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_name VARCHAR(255) NOT NULL,
    institution VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    report_datetime TIMESTAMP NOT NULL,
    area_description TEXT,
    area_category VARCHAR(100) NOT NULL,
    violation_type VARCHAR(150) NOT NULL,
    violation_level VARCHAR(50) NOT NULL,
    environmental_impact VARCHAR(100) NOT NULL,
    urgency_level VARCHAR(20) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    address TEXT,
    status VARCHAR(50) DEFAULT 'PENDING',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_spatial_planning_institution ON spatial_planning_reports(institution);
CREATE INDEX idx_spatial_planning_area_category ON spatial_planning_reports(area_category);
CREATE INDEX idx_spatial_planning_violation_type ON spatial_planning_reports(violation_type);
CREATE INDEX idx_spatial_planning_violation_level ON spatial_planning_reports(violation_level);
CREATE INDEX idx_spatial_planning_urgency_level ON spatial_planning_reports(urgency_level);
CREATE INDEX idx_spatial_planning_status ON spatial_planning_reports(status);
CREATE INDEX idx_spatial_planning_report_datetime ON spatial_planning_reports(report_datetime);

CREATE TABLE spatial_planning_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES spatial_planning_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    caption VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_spatial_planning_photos_report_id ON spatial_planning_photos(report_id);

-- +goose Down
DROP TABLE IF EXISTS spatial_planning_photos;
DROP TABLE IF EXISTS spatial_planning_reports;