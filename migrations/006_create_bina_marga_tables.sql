-- migrations/006_create_bina_marga_tables.sql
-- +goose Up
CREATE TABLE bina_marga_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_name VARCHAR(255) NOT NULL,
    institution_unit VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    report_datetime TIMESTAMP NOT NULL,
    road_name VARCHAR(255) NOT NULL,
    road_type VARCHAR(50) NOT NULL,
    road_class VARCHAR(50) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    damage_type VARCHAR(100) NOT NULL,
    damage_level VARCHAR(50) NOT NULL,
    damaged_length DECIMAL(10, 2) DEFAULT 0,
    damaged_width DECIMAL(10, 2) DEFAULT 0,
    damaged_area DECIMAL(10, 2) DEFAULT 0,
    traffic_impact VARCHAR(50) NOT NULL,
    urgency_level VARCHAR(50) NOT NULL,
    cause_of_damage TEXT,
    status VARCHAR(50) DEFAULT 'PENDING',
    notes TEXT,
    handling_recommendation TEXT,
    estimated_budget DECIMAL(15, 2) DEFAULT 0,
    estimated_repair_time INTEGER DEFAULT 0,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_bina_marga_institution ON bina_marga_reports(institution_unit);
CREATE INDEX idx_bina_marga_road_type ON bina_marga_reports(road_type);
CREATE INDEX idx_bina_marga_road_class ON bina_marga_reports(road_class);
CREATE INDEX idx_bina_marga_road_name ON bina_marga_reports(road_name);
CREATE INDEX idx_bina_marga_damage_type ON bina_marga_reports(damage_type);
CREATE INDEX idx_bina_marga_damage_level ON bina_marga_reports(damage_level);
CREATE INDEX idx_bina_marga_urgency ON bina_marga_reports(urgency_level);
CREATE INDEX idx_bina_marga_traffic_impact ON bina_marga_reports(traffic_impact);
CREATE INDEX idx_bina_marga_status ON bina_marga_reports(status);
CREATE INDEX idx_bina_marga_created_by ON bina_marga_reports(created_by);
CREATE INDEX idx_bina_marga_report_datetime ON bina_marga_reports(report_datetime);

-- Geospatial index for location-based queries
CREATE INDEX idx_bina_marga_location ON bina_marga_reports(latitude, longitude);

-- Priority calculation index for faster sorting
CREATE INDEX idx_bina_marga_priority ON bina_marga_reports(urgency_level, damage_level, traffic_impact, road_class);

-- Composite index for common filter combinations
CREATE INDEX idx_bina_marga_filters ON bina_marga_reports(status, urgency_level, damage_level, created_at);

CREATE TABLE bina_marga_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES bina_marga_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_angle VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bina_marga_photos_report_id ON bina_marga_photos(report_id);

-- Create function to automatically calculate damaged area
CREATE OR REPLACE FUNCTION calculate_damaged_area()
RETURNS TRIGGER AS $$
BEGIN
    NEW.damaged_area = NEW.damaged_length * NEW.damaged_width;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically calculate damaged area on insert/update
CREATE TRIGGER trigger_calculate_damaged_area
    BEFORE INSERT OR UPDATE OF damaged_length, damaged_width
    ON bina_marga_reports
    FOR EACH ROW
    EXECUTE FUNCTION calculate_damaged_area();

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER trigger_update_bina_marga_updated_at
    BEFORE UPDATE ON bina_marga_reports
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create view for priority reports (commonly used query)
CREATE VIEW bina_marga_priority_view AS
SELECT 
    *,
    (CASE WHEN urgency_level = 'DARURAT' THEN 100
          WHEN urgency_level = 'TINGGI' THEN 75
          WHEN urgency_level = 'SEDANG' THEN 50
          ELSE 25 END +
     CASE WHEN damage_level = 'BERAT' THEN 50
          WHEN damage_level = 'SEDANG' THEN 30
          ELSE 15 END +
     CASE WHEN road_class = 'ARTERI' THEN 40
          WHEN road_class = 'KOLEKTOR' THEN 30
          WHEN road_class = 'LOKAL' THEN 20
          ELSE 10 END +
     CASE WHEN traffic_impact = 'TERPUTUS' THEN 60
          WHEN traffic_impact = 'SANGAT_TERGANGGU' THEN 40
          WHEN traffic_impact = 'TERGANGGU' THEN 20
          ELSE 5 END +
     CASE WHEN damaged_area > 100 THEN 25
          WHEN damaged_area > 50 THEN 15
          ELSE 0 END) as priority_score
FROM bina_marga_reports
WHERE status NOT IN ('COMPLETED', 'REJECTED');

-- Create view for statistics (commonly used aggregations)
CREATE VIEW bina_marga_stats_view AS
SELECT 
    COUNT(*) as total_reports,
    COUNT(CASE WHEN urgency_level = 'DARURAT' AND status NOT IN ('COMPLETED', 'REJECTED') THEN 1 END) as emergency_reports,
    COUNT(CASE WHEN traffic_impact = 'TERPUTUS' AND status NOT IN ('COMPLETED', 'REJECTED') THEN 1 END) as blocked_roads,
    COALESCE(SUM(damaged_area), 0) as total_damaged_area,
    COALESCE(SUM(damaged_length), 0) as total_damaged_length,
    COALESCE(SUM(CASE WHEN status NOT IN ('COMPLETED', 'REJECTED') THEN estimated_budget ELSE 0 END), 0) as total_pending_budget,
    COALESCE(AVG(CASE WHEN estimated_repair_time > 0 THEN estimated_repair_time END), 0) as avg_repair_time
FROM bina_marga_reports;

-- +goose Down
DROP VIEW IF EXISTS bina_marga_stats_view;
DROP VIEW IF EXISTS bina_marga_priority_view;
DROP TRIGGER IF EXISTS trigger_update_bina_marga_updated_at ON bina_marga_reports;
DROP TRIGGER IF EXISTS trigger_calculate_damaged_area ON bina_marga_reports;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS calculate_damaged_area();
DROP TABLE IF EXISTS bina_marga_photos;
DROP TABLE IF EXISTS bina_marga_reports;