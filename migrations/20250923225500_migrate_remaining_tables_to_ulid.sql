-- migrations/20250923225500_migrate_remaining_tables_to_ulid.sql
-- +goose Up

-- Recreate water_resources_reports table with ULID
CREATE TABLE water_resources_reports (
    id VARCHAR(26) PRIMARY KEY,
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
    estimated_depth DECIMAL(10, 2) DEFAULT 0,
    estimated_area DECIMAL(10, 2) DEFAULT 0,
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

-- Create indexes for water_resources_reports
CREATE INDEX idx_water_resources_institution ON water_resources_reports(institution_unit);
CREATE INDEX idx_water_resources_irrigation_type ON water_resources_reports(irrigation_type);
CREATE INDEX idx_water_resources_damage_type ON water_resources_reports(damage_type);
CREATE INDEX idx_water_resources_damage_level ON water_resources_reports(damage_level);
CREATE INDEX idx_water_resources_urgency ON water_resources_reports(urgency_category);
CREATE INDEX idx_water_resources_status ON water_resources_reports(status);
CREATE INDEX idx_water_resources_report_datetime ON water_resources_reports(report_datetime);
CREATE INDEX idx_water_resources_irrigation_area ON water_resources_reports(irrigation_area_name);
CREATE INDEX idx_water_resources_priority ON water_resources_reports(urgency_category, damage_level, affected_rice_field_area, affected_farmers_count);

-- Recreate water_resources_photos table with ULID
CREATE TABLE water_resources_photos (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL REFERENCES water_resources_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_angle VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_water_resources_photos_report_id ON water_resources_photos(report_id);

-- Recreate bina_marga_reports table with ULID
CREATE TABLE bina_marga_reports (
    id VARCHAR(26) PRIMARY KEY,
    reporter_name VARCHAR(255) NOT NULL,
    institution_unit VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    report_datetime TIMESTAMP NOT NULL,
    district VARCHAR(255) NOT NULL,
    road_name VARCHAR(255) NOT NULL,
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
    bridge_section VARCHAR(255),
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

-- Create indexes for bina_marga_reports
CREATE INDEX idx_bina_marga_institution ON bina_marga_reports(institution_unit);
CREATE INDEX idx_bina_marga_road_name ON bina_marga_reports(road_name);
CREATE INDEX idx_bina_marga_district ON bina_marga_reports(district);
CREATE INDEX idx_bina_marga_pavement_type ON bina_marga_reports(pavement_type);
CREATE INDEX idx_bina_marga_damage_type ON bina_marga_reports(damage_type);
CREATE INDEX idx_bina_marga_damage_level ON bina_marga_reports(damage_level);
CREATE INDEX idx_bina_marga_urgency ON bina_marga_reports(urgency_level);
CREATE INDEX idx_bina_marga_traffic_impact ON bina_marga_reports(traffic_impact);
CREATE INDEX idx_bina_marga_traffic_condition ON bina_marga_reports(traffic_condition);
CREATE INDEX idx_bina_marga_bridge_name ON bina_marga_reports(bridge_name);
CREATE INDEX idx_bina_marga_bridge_section ON bina_marga_reports(bridge_section);
CREATE INDEX idx_bina_marga_status ON bina_marga_reports(status);
CREATE INDEX idx_bina_marga_report_datetime ON bina_marga_reports(report_datetime);
CREATE INDEX idx_bina_marga_location ON bina_marga_reports(latitude, longitude);
CREATE INDEX idx_bina_marga_priority ON bina_marga_reports(urgency_level, damage_level, traffic_impact);
CREATE INDEX idx_bina_marga_filters ON bina_marga_reports(status, urgency_level, damage_level, created_at);
CREATE INDEX idx_bina_marga_bridges ON bina_marga_reports(bridge_name)
  WHERE bridge_name IS NOT NULL AND bridge_name <> '';
CREATE INDEX idx_bina_marga_emergency ON bina_marga_reports(urgency_level, status)
  WHERE urgency_level = 'DARURAT' AND status NOT IN ('COMPLETED', 'REJECTED');
CREATE INDEX idx_bina_marga_blocked ON bina_marga_reports(traffic_impact, traffic_condition, status)
  WHERE (traffic_impact = 'TERPUTUS' OR traffic_condition = 'TIDAK_BISA_DILALUI_PUTUS')
    AND status NOT IN ('COMPLETED', 'REJECTED');

-- Recreate bina_marga_photos table with ULID
CREATE TABLE bina_marga_photos (
    id VARCHAR(26) PRIMARY KEY,
    report_id VARCHAR(26) NOT NULL REFERENCES bina_marga_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_angle VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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