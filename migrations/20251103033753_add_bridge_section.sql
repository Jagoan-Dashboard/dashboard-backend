-- +goose Up
ALTER TABLE bina_marga_reports 
ADD COLUMN bridge_section VARCHAR(255);

-- +goose Down
ALTER TABLE bina_marga_reports 
DROP COLUMN bridge_section;