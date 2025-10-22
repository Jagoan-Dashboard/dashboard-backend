-- migrations/006_create_bina_marga_tables_updated.sql

-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE bina_marga_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_name VARCHAR(255) NOT NULL,
    institution_unit VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    report_datetime TIMESTAMP NOT NULL,
    road_name VARCHAR(255) NOT NULL,
    road_type VARCHAR(50) NOT NULL,
    road_class VARCHAR(50) NOT NULL,
    segment_length DECIMAL(10, 2) DEFAULT 0,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    pavement_type VARCHAR(50) NOT NULL,
    damage_type VARCHAR(100) NOT NULL,
    damage_level VARCHAR(50) NOT NULL,
    damaged_length DECIMAL(10, 2) DEFAULT 0,
    damaged_width DECIMAL(10, 2) DEFAULT 0,
    damaged_area DECIMAL(10, 2) DEFAULT 0,
    total_damaged_area DECIMAL(10, 2) DEFAULT 0,
    bridge_name VARCHAR(255),
    bridge_structure_type VARCHAR(50),
    bridge_damage_type VARCHAR(100),
    bridge_damage_level VARCHAR(50),
    traffic_condition VARCHAR(50) NOT NULL,
    traffic_impact VARCHAR(50) NOT NULL,
    daily_traffic_volume INTEGER DEFAULT 0,
    urgency_level VARCHAR(50) NOT NULL,
    cause_of_damage TEXT,
    status VARCHAR(50) DEFAULT 'PENDING',
    notes TEXT,
    handling_recommendation TEXT,
    estimated_budget DECIMAL(15, 2) DEFAULT 0,
    estimated_repair_time INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- indexes (tetap sama punyamu) ...
CREATE INDEX idx_bina_marga_institution ON bina_marga_reports(institution_unit);
CREATE INDEX idx_bina_marga_road_type ON bina_marga_reports(road_type);
CREATE INDEX idx_bina_marga_road_class ON bina_marga_reports(road_class);
CREATE INDEX idx_bina_marga_road_name ON bina_marga_reports(road_name);
CREATE INDEX idx_bina_marga_pavement_type ON bina_marga_reports(pavement_type);
CREATE INDEX idx_bina_marga_damage_type ON bina_marga_reports(damage_type);
CREATE INDEX idx_bina_marga_damage_level ON bina_marga_reports(damage_level);
CREATE INDEX idx_bina_marga_urgency ON bina_marga_reports(urgency_level);
CREATE INDEX idx_bina_marga_traffic_impact ON bina_marga_reports(traffic_impact);
CREATE INDEX idx_bina_marga_traffic_condition ON bina_marga_reports(traffic_condition);
CREATE INDEX idx_bina_marga_bridge_name ON bina_marga_reports(bridge_name);
CREATE INDEX idx_bina_marga_status ON bina_marga_reports(status);
CREATE INDEX idx_bina_marga_report_datetime ON bina_marga_reports(report_datetime);
CREATE INDEX idx_bina_marga_location ON bina_marga_reports(latitude, longitude);
CREATE INDEX idx_bina_marga_priority ON bina_marga_reports(urgency_level, damage_level, traffic_impact, road_class);
CREATE INDEX idx_bina_marga_filters ON bina_marga_reports(status, urgency_level, damage_level, created_at);
CREATE INDEX idx_bina_marga_bridges ON bina_marga_reports(bridge_name)
  WHERE bridge_name IS NOT NULL AND bridge_name <> '';
CREATE INDEX idx_bina_marga_emergency ON bina_marga_reports(urgency_level, status)
  WHERE urgency_level = 'DARURAT' AND status NOT IN ('COMPLETED', 'REJECTED');
CREATE INDEX idx_bina_marga_blocked ON bina_marga_reports(traffic_impact, traffic_condition, status)
  WHERE (traffic_impact = 'TERPUTUS' OR traffic_condition = 'TIDAK_BISA_DILALUI_PUTUS')
    AND status NOT IN ('COMPLETED', 'REJECTED');

CREATE TABLE bina_marga_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES bina_marga_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_angle VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_bina_marga_photos_report_id ON bina_marga_photos(report_id);

-- FUNCTION 1: hitung damaged_area  ✅ wrap StatementBegin/End
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION calculate_bina_marga_damaged_area()
RETURNS TRIGGER AS $$
BEGIN
    NEW.damaged_area := COALESCE(NEW.damaged_length, 0) * COALESCE(NEW.damaged_width, 0);

    IF NEW.total_damaged_area IS NULL OR NEW.total_damaged_area = 0 THEN
        NEW.total_damaged_area := NEW.damaged_area;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trigger_calculate_bina_marga_damaged_area
BEFORE INSERT OR UPDATE OF damaged_length, damaged_width, total_damaged_area
ON bina_marga_reports
FOR EACH ROW
EXECUTE FUNCTION calculate_bina_marga_damaged_area();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_bina_marga_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trigger_update_bina_marga_updated_at
BEFORE UPDATE ON bina_marga_reports
FOR EACH ROW
EXECUTE FUNCTION update_bina_marga_updated_at_column();

CREATE VIEW bina_marga_stats_view AS
SELECT 
    COUNT(*) as total_reports,
    COUNT(*) FILTER (WHERE urgency_level = 'DARURAT' AND status NOT IN ('COMPLETED', 'REJECTED')) as emergency_reports,
    COUNT(*) FILTER (WHERE (traffic_impact = 'TERPUTUS' OR traffic_condition = 'TIDAK_BISA_DILALUI_PUTUS') AND status NOT IN ('COMPLETED', 'REJECTED')) as blocked_roads,
    COUNT(*) FILTER (WHERE bridge_name IS NOT NULL AND bridge_name <> '') as bridge_reports,
    COALESCE(SUM(COALESCE(total_damaged_area, damaged_area)), 0) as total_damaged_area,
    COALESCE(SUM(damaged_length), 0) as total_damaged_length,
    COALESCE(SUM(CASE WHEN status NOT IN ('COMPLETED', 'REJECTED') THEN estimated_budget ELSE 0 END), 0) as total_pending_budget,
    COALESCE(AVG(NULLIF(estimated_repair_time, 0)), 0) as avg_repair_time
FROM bina_marga_reports;

CREATE VIEW bina_marga_road_analysis_view AS
SELECT 
    road_type,
    road_class,
    pavement_type,
    COUNT(*) as report_count,
    COUNT(*) FILTER (WHERE status NOT IN ('COMPLETED', 'REJECTED')) as pending_count,
    AVG(COALESCE(total_damaged_area, damaged_area)) as avg_damaged_area,
    SUM(COALESCE(total_damaged_area, damaged_area)) as total_damaged_area,
    AVG(estimated_budget) as avg_estimated_budget,
    AVG(estimated_repair_time) as avg_repair_time
FROM bina_marga_reports
GROUP BY road_type, road_class, pavement_type;

CREATE VIEW bina_marga_bridge_analysis_view AS
SELECT 
    bridge_structure_type,
    bridge_damage_type,
    bridge_damage_level,
    COUNT(*) as bridge_report_count,
    AVG(estimated_budget) as avg_bridge_repair_cost,
    AVG(estimated_repair_time) as avg_bridge_repair_time
FROM bina_marga_reports
WHERE bridge_name IS NOT NULL AND bridge_name <> ''
GROUP BY bridge_structure_type, bridge_damage_type, bridge_damage_level;

-- +goose Down
DROP VIEW IF EXISTS bina_marga_bridge_analysis_view;
DROP VIEW IF EXISTS bina_marga_road_analysis_view;
DROP VIEW IF EXISTS bina_marga_stats_view;
DROP VIEW IF EXISTS bina_marga_priority_view;
DROP TRIGGER IF EXISTS trigger_update_bina_marga_updated_at ON bina_marga_reports;
DROP TRIGGER IF EXISTS trigger_calculate_bina_marga_damaged_area ON bina_marga_reports;
-- +goose StatementBegin
DROP FUNCTION IF EXISTS update_bina_marga_updated_at_column();
DROP FUNCTION IF EXISTS calculate_bina_marga_damaged_area();
-- +goose StatementEnd
DROP TABLE IF EXISTS bina_marga_photos;
DROP TABLE IF EXISTS bina_marga_reports;
