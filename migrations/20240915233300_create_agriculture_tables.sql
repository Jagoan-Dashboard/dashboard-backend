-- migrations/007_create_agriculture_tables.sql
-- +goose Up
CREATE TABLE agriculture_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    extension_officer VARCHAR(255) NOT NULL,
    visit_date DATE NOT NULL,
    farmer_name VARCHAR(255) NOT NULL,
    farmer_group VARCHAR(255),
    farmer_group_type VARCHAR(50),
    village VARCHAR(255) NOT NULL,
    district VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Food Crops Section
    food_commodity VARCHAR(50),
    food_land_status VARCHAR(50),
    food_land_area DECIMAL(10, 2) DEFAULT 0,
    food_growth_phase VARCHAR(50),
    food_plant_age INTEGER DEFAULT 0,
    food_planting_date DATE,
    food_harvest_date DATE,
    food_delay_reason VARCHAR(100),
    food_technology VARCHAR(100),
    
    -- Horticulture Section
    horti_commodity VARCHAR(50),
    horti_sub_commodity VARCHAR(100),
    horti_land_status VARCHAR(50),
    horti_land_area DECIMAL(10, 2) DEFAULT 0,
    horti_growth_phase VARCHAR(50),
    horti_plant_age INTEGER DEFAULT 0,
    horti_planting_date DATE,
    horti_harvest_date DATE,
    horti_delay_reason VARCHAR(100),
    horti_technology VARCHAR(100),
    post_harvest_problems VARCHAR(100),
    
    -- Plantation Section
    plantation_commodity VARCHAR(50),
    plantation_land_status VARCHAR(50),
    plantation_land_area DECIMAL(10, 2) DEFAULT 0,
    plantation_growth_phase VARCHAR(50),
    plantation_plant_age INTEGER DEFAULT 0,
    plantation_planting_date DATE,
    plantation_harvest_date DATE,
    plantation_delay_reason VARCHAR(100),
    plantation_technology VARCHAR(100),
    production_problems VARCHAR(100),
    
    -- Pest and Disease
    has_pest_disease BOOLEAN DEFAULT FALSE,
    pest_disease_type VARCHAR(100),
    pest_disease_commodity VARCHAR(50),
    affected_area VARCHAR(50),
    control_action VARCHAR(100),
    
    -- Weather and Environmental
    weather_condition VARCHAR(50) NOT NULL,
    weather_impact VARCHAR(50) NOT NULL,
    main_constraint VARCHAR(50) NOT NULL,
    
    -- Farmer Needs and Aspirations
    farmer_hope VARCHAR(100) NOT NULL,
    training_needed VARCHAR(100) NOT NULL,
    urgent_needs VARCHAR(100) NOT NULL,
    water_access VARCHAR(50) NOT NULL,
    suggestions TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_agriculture_extension_officer ON agriculture_reports(extension_officer);
CREATE INDEX idx_agriculture_village ON agriculture_reports(village);
CREATE INDEX idx_agriculture_district ON agriculture_reports(district);
CREATE INDEX idx_agriculture_farmer_name ON agriculture_reports(farmer_name);
CREATE INDEX idx_agriculture_visit_date ON agriculture_reports(visit_date);
CREATE INDEX idx_agriculture_food_commodity ON agriculture_reports(food_commodity);
CREATE INDEX idx_agriculture_horti_commodity ON agriculture_reports(horti_commodity);
CREATE INDEX idx_agriculture_plantation_commodity ON agriculture_reports(plantation_commodity);
CREATE INDEX idx_agriculture_has_pest_disease ON agriculture_reports(has_pest_disease);
CREATE INDEX idx_agriculture_main_constraint ON agriculture_reports(main_constraint);

-- Geospatial index for location-based queries
CREATE INDEX idx_agriculture_location ON agriculture_reports(latitude, longitude);

CREATE TABLE agriculture_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES agriculture_reports(id) ON DELETE CASCADE,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50),
    caption VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agriculture_photos_report_id ON agriculture_photos(report_id);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER trigger_update_agriculture_updated_at
    BEFORE UPDATE ON agriculture_reports
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS agriculture_photos;
DROP TABLE IF EXISTS agriculture_reports;