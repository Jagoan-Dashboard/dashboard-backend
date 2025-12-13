-- +goose Up
-- +goose StatementBegin
-- Migration: Add additional columns to rice_fields table for lahan_pengairan import

-- Add Year column
ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS year INT;

-- Add Sawah columns (if not exist)
ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS total_rice_field_area DECIMAL(15,2) DEFAULT 0;

-- Add Non-Sawah columns
ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS dryfield_area DECIMAL(15,2) DEFAULT 0;

ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS shifting_cultivation_area DECIMAL(15,2) DEFAULT 0;

ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS unused_land_area DECIMAL(15,2) DEFAULT 0;

ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS total_non_rice_field_area DECIMAL(15,2) DEFAULT 0;

-- Add Total column
ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS total_land_area DECIMAL(15,2) DEFAULT 0;

-- Add metadata column
ALTER TABLE rice_fields 
ADD COLUMN IF NOT EXISTS data_source VARCHAR(20) DEFAULT 'manual';

-- Add comments for documentation
COMMENT ON COLUMN rice_fields.irrigated_rice_fields IS 'Luas Sawah Irigasi (ha)';
COMMENT ON COLUMN rice_fields.rainfed_rice_fields IS 'Luas Sawah Tadah Hujan (ha)';
COMMENT ON COLUMN rice_fields.total_rice_field_area IS 'Luas Lahan Sawah (ha)';
COMMENT ON COLUMN rice_fields.dryfield_area IS 'Luas Lahan Tegal/Kebun (ha)';
COMMENT ON COLUMN rice_fields.shifting_cultivation_area IS 'Luas Lahan Ladang/Huma (ha)';
COMMENT ON COLUMN rice_fields.unused_land_area IS 'Luas Lahan yang Sementara Tidak Diusahakan (ha)';
COMMENT ON COLUMN rice_fields.total_non_rice_field_area IS 'Luas Lahan Bukan Sawah (ha)';
COMMENT ON COLUMN rice_fields.total_land_area IS 'Total Luas Lahan (ha)';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_rice_fields_year ON rice_fields(year);
CREATE INDEX IF NOT EXISTS idx_rice_fields_district_year ON rice_fields(district, year);
CREATE INDEX IF NOT EXISTS idx_rice_fields_data_source ON rice_fields(data_source);

-- Update existing records to calculate totals
UPDATE rice_fields 
SET 
    year = EXTRACT(YEAR FROM date),
    total_rice_field_area = COALESCE(irrigated_rice_fields, 0) + COALESCE(rainfed_rice_fields, 0),
    total_non_rice_field_area = COALESCE(dryfield_area, 0) + COALESCE(shifting_cultivation_area, 0) + COALESCE(unused_land_area, 0),
    total_land_area = COALESCE(irrigated_rice_fields, 0) + COALESCE(rainfed_rice_fields, 0) + 
                      COALESCE(dryfield_area, 0) + COALESCE(shifting_cultivation_area, 0) + COALESCE(unused_land_area, 0)
WHERE year IS NULL OR total_land_area = 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove indexes
DROP INDEX IF EXISTS idx_rice_fields_data_source;
DROP INDEX IF EXISTS idx_rice_fields_district_year;
DROP INDEX IF EXISTS idx_rice_fields_year;

-- Remove columns (in reverse order)
ALTER TABLE rice_fields DROP COLUMN IF EXISTS data_source;
ALTER TABLE rice_fields DROP COLUMN IF EXISTS total_land_area;
ALTER TABLE rice_fields DROP COLUMN IF EXISTS total_non_rice_field_area;
ALTER TABLE rice_fields DROP COLUMN IF EXISTS unused_land_area;
ALTER TABLE rice_fields DROP COLUMN IF EXISTS shifting_cultivation_area;
ALTER TABLE rice_fields DROP COLUMN IF EXISTS dryfield_area;
ALTER TABLE rice_fields DROP COLUMN IF EXISTS total_rice_field_area;
ALTER TABLE rice_fields DROP COLUMN IF EXISTS year;
-- +goose StatementEnd