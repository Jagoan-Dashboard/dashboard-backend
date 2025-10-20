-- +goose Up
CREATE TABLE indikator_demografi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    indikator VARCHAR(150) NOT NULL,
    tahun INTEGER NOT NULL,
    nilai DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_demografi_indikator_tahun UNIQUE(indikator, tahun)
);

CREATE INDEX idx_indikator_demografi_indikator ON indikator_demografi(indikator);
CREATE INDEX idx_indikator_demografi_tahun ON indikator_demografi(tahun);
CREATE INDEX idx_indikator_demografi_indikator_tahun ON indikator_demografi(indikator, tahun);

COMMENT ON TABLE indikator_demografi IS 'Tabel untuk menyimpan data indikator demografi per tahun';

-- Seeder data
INSERT INTO indikator_demografi (id, indikator, tahun, nilai) VALUES
-- Kepadatan Penduduk
(gen_random_uuid(), 'Kepadatan Penduduk', 2020, 624),
(gen_random_uuid(), 'Kepadatan Penduduk', 2021, 626),
(gen_random_uuid(), 'Kepadatan Penduduk', 2022, 629),
(gen_random_uuid(), 'Kepadatan Penduduk', 2023, 631),
(gen_random_uuid(), 'Kepadatan Penduduk', 2024, 634),

-- Rasio Ketergantungan
(gen_random_uuid(), 'Rasio Ketergantungan', 2020, 44.33),
(gen_random_uuid(), 'Rasio Ketergantungan', 2021, 44.78),
(gen_random_uuid(), 'Rasio Ketergantungan', 2022, 45.50),
(gen_random_uuid(), 'Rasio Ketergantungan', 2023, 46.23),
(gen_random_uuid(), 'Rasio Ketergantungan', 2024, 46.96),

-- Jumlah penduduk produktif
(gen_random_uuid(), 'Jumlah penduduk produktif (Usia 15–64 tahun)', 2020, 603147),
(gen_random_uuid(), 'Jumlah penduduk produktif (Usia 15–64 tahun)', 2021, 603093),
(gen_random_uuid(), 'Jumlah penduduk produktif (Usia 15–64 tahun)', 2022, 602781),
(gen_random_uuid(), 'Jumlah penduduk produktif (Usia 15–64 tahun)', 2023, 602263),
(gen_random_uuid(), 'Jumlah penduduk produktif (Usia 15–64 tahun)', 2024, 601571),

-- Jumlah penduduk non produktif
(gen_random_uuid(), 'Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)', 2020, 267384),
(gen_random_uuid(), 'Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)', 2021, 270058),
(gen_random_uuid(), 'Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)', 2022, 274288),
(gen_random_uuid(), 'Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)', 2023, 278450),
(gen_random_uuid(), 'Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)', 2024, 282515);

-- +goose Down
DROP INDEX IF EXISTS idx_indikator_demografi_indikator_tahun;
DROP INDEX IF EXISTS idx_indikator_demografi_tahun;
DROP INDEX IF EXISTS idx_indikator_demografi_indikator;
DROP TABLE IF EXISTS indikator_demografi;