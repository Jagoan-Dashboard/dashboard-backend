-- +goose Up
ALTER TABLE bina_marga_reports 
ALTER COLUMN reporter_name DROP NOT NULL,
ALTER COLUMN road_type DROP NOT NULL,
ALTER COLUMN traffic_impact DROP NOT NULL,
ALTER COLUMN report_date_time DROP NOT NULL,
ALTER COLUMN road_class DROP NOT NULL,
ALTER COLUMN created_by DROP NOT NULL,
ALTER COLUMN report_datetime DROP NOT NULL;

-- +goose Down
ALTER TABLE bina_marga_reports 
ALTER COLUMN reporter_name SET NOT NULL,
ALTER COLUMN road_type SET NOT NULL,
ALTER COLUMN traffic_impact SET NOT NULL,
ALTER COLUMN report_date_time SET NOT NULL,
ALTER COLUMN road_class SET NOT NULL,
ALTER COLUMN created_by SET NOT NULL,
ALTER COLUMN report_datetime SET NOT NULL;