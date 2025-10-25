-- migrations/20251023120000_add_district_and_update_enums.sql

-- +goose Up

ALTER TABLE bina_marga_reports 
ADD COLUMN IF NOT EXISTS district VARCHAR(255);

COMMENT ON COLUMN bina_marga_reports.district IS 'Kecamatan/Lokasi Jalan';

-- Create index for district filtering
CREATE INDEX IF NOT EXISTS idx_bina_marga_district ON bina_marga_reports(district);


ALTER TABLE bina_marga_reports 
ALTER COLUMN road_type DROP NOT NULL;

ALTER TABLE bina_marga_reports 
ALTER COLUMN road_class DROP NOT NULL;


ALTER TABLE bina_marga_reports 
ALTER COLUMN traffic_impact DROP NOT NULL;


UPDATE bina_marga_reports 
SET district = 'UNKNOWN' 
WHERE district IS NULL OR district = '';

DROP INDEX IF EXISTS idx_bina_marga_district;
ALTER TABLE bina_marga_reports DROP COLUMN IF EXISTS district;

ALTER TABLE bina_marga_reports 
ALTER COLUMN road_type SET NOT NULL;

ALTER TABLE bina_marga_reports 
ALTER COLUMN road_class SET NOT NULL;

ALTER TABLE bina_marga_reports 
ALTER COLUMN traffic_impact SET NOT NULL;

