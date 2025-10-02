-- migrations/005_create_water_resources_tables.sql
-- +goose Up
CREATE TABLE water_resources_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_name VARCHAR(255) NOT NULL,
    institution_unit VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    report_datetime TIMESTAMP NOT NULL,
    irrigation_area_name VARCHAR(255),
    irrigation_type VARCHAR(50) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    damage_type VARCHAR(100) NOT NULL,
    damage_level VARCHAR(50) NOT NULL,
    estimated_length DECIMAL(10, 2) DEFAULT 0,
    estimated_width DECIMAL(10, 2) DEFAULT 0,
    estimated_volume DECIMAL(10, 2) DEFAULT 0,
    affected_rice_field_area DECIMAL(10, 2) DEFAULT 0,
    affected_farmers_count INTEGER DEFAULT 0,
    urgency_category VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING',
    notes TEXT,
    handling_recommendation TEXT,
    estimated_budget DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_water_resources_institution ON water_resources_reports(institution_unit);
CREATE INDEX idx_water_resources_irrigation_type ON water_resources_reports(irrigation_type);
CREATE INDEX idx_water_resources_damage_type ON water_resources_reports(damage_type);
CREATE INDEX idx_water_resources_damage_level ON water_resources_reports(damage_level);
CREATE INDEX idx_water_resources_urgency ON water_resources_reports(urgency_category);
CREATE INDEX idx_water_resources_status ON water_resources_reports(status);
CREATE INDEX idx_water_resources_report_datetime ON water_resources_reports(report_datetime);
CREATE INDEX idx_water_resources_irrigation_area ON water_resources_reports(irrigation_area_name);

-- Priority calculation index for faster sorting
CREATE INDEX idx_water_resources_priority ON water_resources_reports(urgency_category, damage_level, affected_rice_field_area, affected_farmers_count);

CREATE TABLE water_resources_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES water_resources_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_angle VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_water_resources_photos_report_id ON water_resources_photos(report_id);

-- +goose Down
DROP TABLE IF EXISTS water_resources_photos;
DROP TABLE IF EXISTS water_resources_reports;