-- +goose Up
ALTER TABLE spatial_planning_reports 
ALTER COLUMN reporter_name DROP NOT NULL,
ALTER COLUMN district DROP NOT NULL,
ALTER COLUMN report_date_time DROP NOT NULL,
ALTER COLUMN created_by DROP NOT NULL,
ALTER COLUMN report_datetime DROP NOT NULL;

-- +goose Down
ALTER TABLE spatial_planning_reports 
ALTER COLUMN reporter_name SET NOT NULL,
ALTER COLUMN district SET NOT NULL,
ALTER COLUMN report_date_time SET NOT NULL,
ALTER COLUMN created_by SET NOT NULL,
ALTER COLUMN report_datetime SET NOT NULL;