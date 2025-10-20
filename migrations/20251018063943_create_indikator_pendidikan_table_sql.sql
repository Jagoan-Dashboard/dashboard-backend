-- +goose Up
CREATE TABLE indikator_pendidikan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
(gen_random_uuid(), 'Rata-rata Lama Sekolah', 2020, 7.06),
(gen_random_uuid(), 'Rata-rata Lama Sekolah', 2021, 7.26),
(gen_random_uuid(), 'Rata-rata Lama Sekolah', 2022, 7.59),
(gen_random_uuid(), 'Rata-rata Lama Sekolah', 2023, 7.78),
(gen_random_uuid(), 'Rata-rata Lama Sekolah', 2024, 7.84),

-- Harapan Lama Sekolah
(gen_random_uuid(), 'Harapan Lama Sekolah', 2020, 12.70),
(gen_random_uuid(), 'Harapan Lama Sekolah', 2021, 12.83),
(gen_random_uuid(), 'Harapan Lama Sekolah', 2022, 12.84),
(gen_random_uuid(), 'Harapan Lama Sekolah', 2023, 12.85),
(gen_random_uuid(), 'Harapan Lama Sekolah', 2024, 12.89),

-- Proporsi Penduduk dengan Pendidikan Tinggi
(gen_random_uuid(), 'Proporsi dengan Pendidikan Tinggi', 2020, 5.23),
(gen_random_uuid(), 'Proporsi dengan Pendidikan Tinggi', 2021, 6.79),
(gen_random_uuid(), 'Proporsi dengan Pendidikan Tinggi', 2022, 6.87),
(gen_random_uuid(), 'Proporsi dengan Pendidikan Tinggi', 2023, 7.42),
(gen_random_uuid(), 'Proporsi dengan Pendidikan Tinggi', 2024, 7.50);

-- +goose Down
DROP INDEX IF EXISTS idx_indikator_pendidikan_indikator_tahun;
DROP INDEX IF EXISTS idx_indikator_pendidikan_tahun;
DROP INDEX IF EXISTS idx_indikator_pendidikan_indikator;
DROP TABLE IF EXISTS indikator_pendidikan;