-- +goose Up
ALTER TABLE water_resources_reports 
ADD COLUMN estimated_depth NUMERIC;

-- +goose Down
ALTER TABLE water_resources_reports 
DROP COLUMN estimated_depth;