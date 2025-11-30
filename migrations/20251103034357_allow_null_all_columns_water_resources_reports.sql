-- +goose Up
ALTER TABLE water_resources_reports 
ALTER COLUMN reporter_name DROP NOT NULL,
ALTER COLUMN report_datetime DROP NOT NULL;

-- +goose Down
ALTER TABLE water_resources_reports 
ALTER COLUMN reporter_name SET NOT NULL,
ALTER COLUMN report_datetime SET NOT NULL;