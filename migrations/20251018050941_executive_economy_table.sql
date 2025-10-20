-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE indikator_ekonomi (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    indikator VARCHAR(100) NOT NULL,
    tahun INTEGER NOT NULL,
    nilai DECIMAL(10, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_indikator_tahun UNIQUE(indikator, tahun)
);

CREATE INDEX idx_indikator_ekonomi_indikator ON indikator_ekonomi(indikator);
CREATE INDEX idx_indikator_ekonomi_tahun ON indikator_ekonomi(tahun);
CREATE INDEX idx_indikator_ekonomi_indikator_tahun ON indikator_ekonomi(indikator, tahun);

COMMENT ON TABLE indikator_ekonomi IS 'Tabel untuk menyimpan data indikator ekonomi per tahun';
COMMENT ON COLUMN indikator_ekonomi.id IS 'Primary key UUID';
COMMENT ON COLUMN indikator_ekonomi.indikator IS 'Nama indikator ekonomi (Laju Pertumbuhan Ekonomi, % Pertanian (PDRB), dll)';
COMMENT ON COLUMN indikator_ekonomi.tahun IS 'Tahun data indikator';
COMMENT ON COLUMN indikator_ekonomi.nilai IS 'Nilai indikator dengan 4 digit desimal';

-- Seeder data indikator ekonomi (2020–2024)
INSERT INTO indikator_ekonomi (id, indikator, tahun, nilai) VALUES
-- Laju Pertumbuhan Ekonomi
(uuid_generate_v4(), 'Laju Pertumbuhan Ekonomi', 2020, -1.69),
(uuid_generate_v4(), 'Laju Pertumbuhan Ekonomi', 2021, 2.55),
(uuid_generate_v4(), 'Laju Pertumbuhan Ekonomi', 2022, 3.19),
(uuid_generate_v4(), 'Laju Pertumbuhan Ekonomi', 2023, 4.49),
(uuid_generate_v4(), 'Laju Pertumbuhan Ekonomi', 2024, 4.64),

-- % Pertanian (PDRB)
(uuid_generate_v4(), '% Pertanian (PDRB)', 2020, 35.33),
(uuid_generate_v4(), '% Pertanian (PDRB)', 2021, 33.80),
(uuid_generate_v4(), '% Pertanian (PDRB)', 2022, 32.93),
(uuid_generate_v4(), '% Pertanian (PDRB)', 2023, 32.86),
(uuid_generate_v4(), '% Pertanian (PDRB)', 2024, 32.04),

-- % Pengolahan (PDRB)
(uuid_generate_v4(), '% Pengolahan (PDRB)', 2020, 8.73),
(uuid_generate_v4(), '% Pengolahan (PDRB)', 2021, 9.31),
(uuid_generate_v4(), '% Pengolahan (PDRB)', 2022, 9.72),
(uuid_generate_v4(), '% Pengolahan (PDRB)', 2023, 9.87),
(uuid_generate_v4(), '% Pengolahan (PDRB)', 2024, 10.12),

-- ICOR
(uuid_generate_v4(), 'ICOR', 2020, -10.58),
(uuid_generate_v4(), 'ICOR', 2021, 9.25),
(uuid_generate_v4(), 'ICOR', 2022, 5.72),
(uuid_generate_v4(), 'ICOR', 2023, 6.05),
(uuid_generate_v4(), 'ICOR', 2024, 6.17),

-- ILOR
(uuid_generate_v4(), 'ILOR', 2020, 0.015),
(uuid_generate_v4(), 'ILOR', 2021, 0.005),
(uuid_generate_v4(), 'ILOR', 2022, 0.004),
(uuid_generate_v4(), 'ILOR', 2023, 0.003),
(uuid_generate_v4(), 'ILOR', 2024, 0.009),

-- Inflasi
(uuid_generate_v4(), 'Inflasi', 2020, 1.93),
(uuid_generate_v4(), 'Inflasi', 2021, 1.64),
(uuid_generate_v4(), 'Inflasi', 2022, 5.76),
(uuid_generate_v4(), 'Inflasi', 2023, 2.64),
(uuid_generate_v4(), 'Inflasi', 2024, 1.19);

-- +goose Down
DROP INDEX IF EXISTS idx_indikator_ekonomi_indikator_tahun;
DROP INDEX IF EXISTS idx_indikator_ekonomi_tahun;
DROP INDEX IF EXISTS idx_indikator_ekonomi_indikator;
DROP TABLE IF EXISTS indikator_ekonomi;