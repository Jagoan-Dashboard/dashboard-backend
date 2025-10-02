-- migrations/20250923225400_migrate_to_ulid.sql
-- +goose Up

-- Drop existing tables to recreate with ULID (ensure backup is taken before running this)
-- Note: This migration assumes development environment. For production, data migration script needed.

DROP TABLE IF EXISTS agriculture_photos;
DROP TABLE IF EXISTS agriculture_reports;
DROP TABLE IF EXISTS bina_marga_photos;
DROP TABLE IF EXISTS bina_marga_reports;
DROP TABLE IF EXISTS water_resources_photos;
DROP TABLE IF EXISTS water_resources_reports;
DROP TABLE IF EXISTS spatial_planning_photos;
DROP TABLE IF EXISTS spatial_planning_reports;
DROP TABLE IF EXISTS report_photos;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS users;

-- Recreate users table with ULID
CREATE TABLE users (
    id VARCHAR(26) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'OPERATOR',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for users
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

-- Recreate reports table with ULID
CREATE TABLE reports (
    id VARCHAR(26) PRIMARY KEY,
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
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for reports
CREATE INDEX idx_reports_village ON reports(village);
CREATE INDEX idx_reports_district ON reports(district);
CREATE INDEX idx_reports_building_type ON reports(building_type);
CREATE INDEX idx_reports_report_status ON reports(report_status);
CREATE INDEX idx_reports_created_at ON reports(created_at);

-- Recreate report_photos table with ULID
CREATE TABLE report_photos (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for report_photos
CREATE INDEX idx_report_photos_report_id ON report_photos(report_id);
CREATE INDEX idx_report_photos_photo_type ON report_photos(photo_type);

-- Recreate spatial_planning_reports table with ULID
CREATE TABLE spatial_planning_reports (
    id VARCHAR(26) PRIMARY KEY,
    reporter_name VARCHAR(255) NOT NULL,
    reporter_role VARCHAR(50) NOT NULL,
    village VARCHAR(255) NOT NULL,
    district VARCHAR(255) NOT NULL,
    location_details TEXT,
    violation_type VARCHAR(100) NOT NULL,
    violation_level VARCHAR(50) NOT NULL,
    urgency_level VARCHAR(50) NOT NULL,
    area_category VARCHAR(50) NOT NULL,
    environmental_impact VARCHAR(100),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    total_area DECIMAL(12, 2),
    affected_area DECIMAL(12, 2),
    report_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    priority_score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for spatial_planning_reports
CREATE INDEX idx_spatial_planning_village ON spatial_planning_reports(village);
CREATE INDEX idx_spatial_planning_district ON spatial_planning_reports(district);
CREATE INDEX idx_spatial_planning_violation_type ON spatial_planning_reports(violation_type);
CREATE INDEX idx_spatial_planning_urgency_level ON spatial_planning_reports(urgency_level);
CREATE INDEX idx_spatial_planning_report_status ON spatial_planning_reports(report_status);
CREATE INDEX idx_spatial_planning_priority_score ON spatial_planning_reports(priority_score);

-- Recreate spatial_planning_photos table with ULID
CREATE TABLE spatial_planning_photos (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL REFERENCES spatial_planning_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_spatial_planning_photos_report_id ON spatial_planning_photos(report_id);

-- +goose Down
DROP TABLE IF EXISTS spatial_planning_photos;
DROP TABLE IF EXISTS spatial_planning_reports;
DROP TABLE IF EXISTS report_photos;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS users;