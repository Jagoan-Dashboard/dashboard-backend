-- +goose Up
CREATE TABLE indikator_ketenagakerjaan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    indikator VARCHAR(150) NOT NULL,
    tahun INTEGER NOT NULL,
    nilai DECIMAL(15, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_ketenagakerjaan_indikator_tahun UNIQUE(indikator, tahun)
);

CREATE INDEX idx_indikator_ketenagakerjaan_indikator ON indikator_ketenagakerjaan(indikator);
CREATE INDEX idx_indikator_ketenagakerjaan_tahun ON indikator_ketenagakerjaan(tahun);
CREATE INDEX idx_indikator_ketenagakerjaan_indikator_tahun ON indikator_ketenagakerjaan(indikator, tahun);

COMMENT ON TABLE indikator_ketenagakerjaan IS 'Tabel untuk menyimpan data indikator ketenagakerjaan per tahun';
COMMENT ON COLUMN indikator_ketenagakerjaan.indikator IS 'Nama indikator ketenagakerjaan (TPT, TPAK, dll)';
COMMENT ON COLUMN indikator_ketenagakerjaan.tahun IS 'Tahun data indikator';
COMMENT ON COLUMN indikator_ketenagakerjaan.nilai IS 'Nilai indikator dengan 4 digit desimal';

-- Seeder data
INSERT INTO indikator_ketenagakerjaan (id, indikator, tahun, nilai) VALUES
-- Tingkat Pengangguran Terbuka (TPT)
(gen_random_uuid(), 'TPT', 2020, 5.44),
(gen_random_uuid(), 'TPT', 2021, 4.25),
(gen_random_uuid(), 'TPT', 2022, 2.48),
(gen_random_uuid(), 'TPT', 2023, 2.41),
(gen_random_uuid(), 'TPT', 2024, 2.40),

-- Tingkat Partisipasi Angkatan Kerja (TPAK)
(gen_random_uuid(), 'TPAK', 2020, 72.69),
(gen_random_uuid(), 'TPAK', 2021, 72.88),
(gen_random_uuid(), 'TPAK', 2022, 78.60),
(gen_random_uuid(), 'TPAK', 2023, 69.43),
(gen_random_uuid(), 'TPAK', 2024, 75.73),

-- Tingkat Partisipasi Angkatan Kerja (TPAK) Perempuan
(gen_random_uuid(), 'TPAK Perempuan', 2020, 61.19),
(gen_random_uuid(), 'TPAK Perempuan', 2021, 61.78),
(gen_random_uuid(), 'TPAK Perempuan', 2022, 69.90),
(gen_random_uuid(), 'TPAK Perempuan', 2023, 55.69),
(gen_random_uuid(), 'TPAK Perempuan', 2024, 65.64),

-- Upah Minimum Kabupaten (Rupiah)
(gen_random_uuid(), 'Upah Minimum Kabupaten (Rupiah)', 2020, 1913321.73),
(gen_random_uuid(), 'Upah Minimum Kabupaten (Rupiah)', 2021, 1960510.00),
(gen_random_uuid(), 'Upah Minimum Kabupaten (Rupiah)', 2022, 1962585.99),
(gen_random_uuid(), 'Upah Minimum Kabupaten (Rupiah)', 2023, 2158844.59),
(gen_random_uuid(), 'Upah Minimum Kabupaten (Rupiah)', 2024, 2241054.00);

-- +goose Down
DROP INDEX IF EXISTS idx_indikator_ketenagakerjaan_indikator_tahun;
DROP INDEX IF EXISTS idx_indikator_ketenagakerjaan_tahun;
DROP INDEX IF EXISTS idx_indikator_ketenagakerjaan_indikator;
DROP TABLE IF EXISTS indikator_ketenagakerjaan;