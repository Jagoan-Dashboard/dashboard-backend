-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE indikator_pendidikan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    indikator VARCHAR(150) NOT NULL,
    tahun INTEGER NOT NULL,
    nilai DECIMAL(15, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_pendidikan_indikator_tahun UNIQUE(indikator, tahun)
);

CREATE INDEX idx_indikator_pendidikan_indikator ON indikator_pendidikan(indikator);
CREATE INDEX idx_indikator_pendidikan_tahun ON indikator_pendidikan(tahun);
CREATE INDEX idx_indikator_pendidikan_indikator_tahun ON indikator_pendidikan(indikator, tahun);

COMMENT ON TABLE indikator_pendidikan IS 'Tabel untuk menyimpan data indikator pendidikan per tahun';
COMMENT ON COLUMN indikator_pendidikan.indikator IS 'Nama indikator pendidikan (Rata-rata Lama Sekolah, dll)';
COMMENT ON COLUMN indikator_pendidikan.tahun IS 'Tahun data indikator';
COMMENT ON COLUMN indikator_pendidikan.nilai IS 'Nilai indikator dengan 4 digit desimal';

-- Seeder data
INSERT INTO indikator_pendidikan (id, indikator, tahun, nilai) VALUES
-- Rata-rata Lama Sekolah
(uuid_generate_v4(), 'Rata-rata Lama Sekolah', 2020, 7.06),
(uuid_generate_v4(), 'Rata-rata Lama Sekolah', 2021, 7.26),
(uuid_generate_v4(), 'Rata-rata Lama Sekolah', 2022, 7.59),
(uuid_generate_v4(), 'Rata-rata Lama Sekolah', 2023, 7.78),
(uuid_generate_v4(), 'Rata-rata Lama Sekolah', 2024, 7.84),

-- Harapan Lama Sekolah
(uuid_generate_v4(), 'Harapan Lama Sekolah', 2020, 12.70),
(uuid_generate_v4(), 'Harapan Lama Sekolah', 2021, 12.83),
(uuid_generate_v4(), 'Harapan Lama Sekolah', 2022, 12.84),
(uuid_generate_v4(), 'Harapan Lama Sekolah', 2023, 12.85),
(uuid_generate_v4(), 'Harapan Lama Sekolah', 2024, 12.89),

-- Proporsi Penduduk dengan Pendidikan Tinggi
(uuid_generate_v4(), 'Proporsi dengan Pendidikan Tinggi', 2020, 5.23),
(uuid_generate_v4(), 'Proporsi dengan Pendidikan Tinggi', 2021, 6.79),
(uuid_generate_v4(), 'Proporsi dengan Pendidikan Tinggi', 2022, 6.87),
(uuid_generate_v4(), 'Proporsi dengan Pendidikan Tinggi', 2023, 7.42),
(uuid_generate_v4(), 'Proporsi dengan Pendidikan Tinggi', 2024, 7.50);

-- +goose Down
DROP INDEX IF EXISTS idx_indikator_pendidikan_indikator_tahun;
DROP INDEX IF EXISTS idx_indikator_pendidikan_tahun;
DROP INDEX IF EXISTS idx_indikator_pendidikan_indikator;
DROP TABLE IF EXISTS indikator_pendidikan;