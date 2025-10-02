-- migrations/20250923225500_migrate_remaining_tables_to_ulid.sql
-- +goose Up

-- Recreate water_resources_reports table with ULID
CREATE TABLE water_resources_reports (
    id VARCHAR(26) PRIMARY KEY,
    reporter_name VARCHAR(255) NOT NULL,
    reporter_role VARCHAR(50) NOT NULL,
    village VARCHAR(255) NOT NULL,
    district VARCHAR(255) NOT NULL,
    location_details TEXT,
    water_source_type VARCHAR(50) NOT NULL,
    damage_type VARCHAR(100) NOT NULL,
    damage_severity VARCHAR(50) NOT NULL,
    urgent_level VARCHAR(50) NOT NULL,
    affected_population INTEGER,
    estimated_loss DECIMAL(15, 2),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    water_flow_before DECIMAL(8, 2),
    water_flow_after DECIMAL(8, 2),
    infrastructure_condition VARCHAR(100),
    report_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    priority_score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for water_resources_reports
CREATE INDEX idx_water_resources_village ON water_resources_reports(village);
CREATE INDEX idx_water_resources_district ON water_resources_reports(district);
CREATE INDEX idx_water_resources_damage_type ON water_resources_reports(damage_type);
CREATE INDEX idx_water_resources_urgent_level ON water_resources_reports(urgent_level);
CREATE INDEX idx_water_resources_report_status ON water_resources_reports(report_status);
CREATE INDEX idx_water_resources_priority_score ON water_resources_reports(priority_score);

-- Recreate water_resources_photos table with ULID
CREATE TABLE water_resources_photos (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL REFERENCES water_resources_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_water_resources_photos_report_id ON water_resources_photos(report_id);

-- Recreate bina_marga_reports table with ULID
CREATE TABLE bina_marga_reports (
    id VARCHAR(26) PRIMARY KEY,
    reporter_name VARCHAR(255) NOT NULL,
    reporter_role VARCHAR(50) NOT NULL,
    village VARCHAR(255) NOT NULL,
    district VARCHAR(255) NOT NULL,
    road_name VARCHAR(255) NOT NULL,
    road_type VARCHAR(50) NOT NULL,
    road_status VARCHAR(50) NOT NULL,
    damage_type VARCHAR(100) NOT NULL,
    damage_severity VARCHAR(50) NOT NULL,
    urgent_level VARCHAR(50) NOT NULL,
    road_length DECIMAL(8, 2),
    road_width DECIMAL(6, 2),
    affected_length DECIMAL(8, 2),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    traffic_impact VARCHAR(100),
    estimated_cost DECIMAL(15, 2),
    report_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    priority_score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for bina_marga_reports
CREATE INDEX idx_bina_marga_village ON bina_marga_reports(village);
CREATE INDEX idx_bina_marga_district ON bina_marga_reports(district);
CREATE INDEX idx_bina_marga_road_type ON bina_marga_reports(road_type);
CREATE INDEX idx_bina_marga_damage_type ON bina_marga_reports(damage_type);
CREATE INDEX idx_bina_marga_urgent_level ON bina_marga_reports(urgent_level);
CREATE INDEX idx_bina_marga_report_status ON bina_marga_reports(report_status);
CREATE INDEX idx_bina_marga_priority_score ON bina_marga_reports(priority_score);

-- Recreate bina_marga_photos table with ULID
CREATE TABLE bina_marga_photos (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL REFERENCES bina_marga_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bina_marga_photos_report_id ON bina_marga_photos(report_id);

-- Recreate agriculture_reports table with ULID
CREATE TABLE agriculture_reports (
    id VARCHAR(26) PRIMARY KEY,
    extension_officer VARCHAR(255) NOT NULL,
    visit_date DATE NOT NULL,
    farmer_name VARCHAR(255) NOT NULL,
    farmer_group VARCHAR(255),
    farmer_group_type VARCHAR(50),
    village VARCHAR(255) NOT NULL,
    district VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Food crop fields
    food_commodity VARCHAR(50),
    food_land_status VARCHAR(50),
    food_land_area DECIMAL(10, 3),
    food_growth_phase VARCHAR(50),
    food_plant_age INTEGER,
    food_planting_date DATE,
    food_harvest_date DATE,
    food_delay_reason VARCHAR(100),
    food_technology VARCHAR(100),

    -- Horticulture fields
    horti_commodity VARCHAR(50),
    horti_sub_commodity VARCHAR(100),
    horti_land_status VARCHAR(50),
    horti_land_area DECIMAL(10, 3),
    horti_growth_phase VARCHAR(50),
    horti_plant_age INTEGER,
    horti_planting_date DATE,
    horti_harvest_date DATE,
    horti_delay_reason VARCHAR(100),
    horti_technology VARCHAR(100),
    post_harvest_problems VARCHAR(100),

    -- Plantation fields
    plantation_commodity VARCHAR(50),
    plantation_land_status VARCHAR(50),
    plantation_land_area DECIMAL(10, 3),
    plantation_growth_phase VARCHAR(50),
    plantation_plant_age INTEGER,
    plantation_planting_date DATE,
    plantation_harvest_date DATE,
    plantation_delay_reason VARCHAR(100),
    plantation_technology VARCHAR(100),
    production_problems VARCHAR(100),

    -- Pest and disease
    has_pest_disease BOOLEAN NOT NULL DEFAULT false,
    pest_disease_type VARCHAR(100),
    pest_disease_commodity VARCHAR(50),
    affected_area VARCHAR(50),
    control_action VARCHAR(100),

    -- Environmental and constraints
    weather_condition VARCHAR(50) NOT NULL,
    weather_impact VARCHAR(50) NOT NULL,
    main_constraint VARCHAR(50) NOT NULL,

    -- Farmer needs and hopes
    farmer_hope VARCHAR(100) NOT NULL,
    training_needed VARCHAR(100) NOT NULL,
    urgent_needs VARCHAR(100) NOT NULL,
    water_access VARCHAR(50) NOT NULL,
    suggestions TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for agriculture_reports
CREATE INDEX idx_agriculture_village ON agriculture_reports(village);
CREATE INDEX idx_agriculture_district ON agriculture_reports(district);
CREATE INDEX idx_agriculture_extension_officer ON agriculture_reports(extension_officer);
CREATE INDEX idx_agriculture_visit_date ON agriculture_reports(visit_date);
CREATE INDEX idx_agriculture_farmer_group_type ON agriculture_reports(farmer_group_type);
CREATE INDEX idx_agriculture_food_commodity ON agriculture_reports(food_commodity);
CREATE INDEX idx_agriculture_horti_commodity ON agriculture_reports(horti_commodity);
CREATE INDEX idx_agriculture_plantation_commodity ON agriculture_reports(plantation_commodity);
CREATE INDEX idx_agriculture_has_pest_disease ON agriculture_reports(has_pest_disease);

-- Recreate agriculture_photos table with ULID
CREATE TABLE agriculture_photos (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL REFERENCES agriculture_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agriculture_photos_report_id ON agriculture_photos(report_id);

-- +goose Down
DROP TABLE IF EXISTS agriculture_photos;
DROP TABLE IF EXISTS agriculture_reports;
DROP TABLE IF EXISTS bina_marga_photos;
DROP TABLE IF EXISTS bina_marga_reports;
DROP TABLE IF EXISTS water_resources_photos;
DROP TABLE IF EXISTS water_resources_reports;