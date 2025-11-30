-- +goose Up
ALTER TABLE water_resources_reports 
ADD COLUMN estimated_area NUMERIC;

-- +goose Down
ALTER TABLE water_resources_reports 
DROP COLUMN estimated_area;