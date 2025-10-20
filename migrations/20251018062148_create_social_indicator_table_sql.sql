-- +goose Up
CREATE TABLE indikator_sosial (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    indikator VARCHAR(150) NOT NULL,
    tahun INTEGER NOT NULL,
    nilai DECIMAL(15, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_sosial_indikator_tahun UNIQUE(indikator, tahun)
);

CREATE INDEX idx_indikator_sosial_indikator ON indikator_sosial(indikator);
CREATE INDEX idx_indikator_sosial_tahun ON indikator_sosial(tahun);
CREATE INDEX idx_indikator_sosial_indikator_tahun ON indikator_sosial(indikator, tahun);

COMMENT ON TABLE indikator_sosial IS 'Tabel untuk menyimpan data indikator sosial/kemiskinan per tahun';
COMMENT ON COLUMN indikator_sosial.indikator IS 'Nama indikator sosial (Angka Kemiskinan, IPM, dll)';
COMMENT ON COLUMN indikator_sosial.tahun IS 'Tahun data indikator';
COMMENT ON COLUMN indikator_sosial.nilai IS 'Nilai indikator dengan 4 digit desimal';

-- Seeder data
INSERT INTO indikator_sosial (id, indikator, tahun, nilai) VALUES
-- Angka Kemiskinan (P0)
(gen_random_uuid(), 'Angka Kemiskinan (P0)', 2020, 15.44),
(gen_random_uuid(), 'Angka Kemiskinan (P0)', 2021, 15.57),
(gen_random_uuid(), 'Angka Kemiskinan (P0)', 2022, 14.15),
(gen_random_uuid(), 'Angka Kemiskinan (P0)', 2023, 14.40),
(gen_random_uuid(), 'Angka Kemiskinan (P0)', 2024, 13.81),

-- Indeks Kedalaman Kemiskinan (P1)
(gen_random_uuid(), 'Indeks Kedalaman Kemiskinan (P1)', 2020, 2.37),
(gen_random_uuid(), 'Indeks Kedalaman Kemiskinan (P1)', 2021, 2.23),
(gen_random_uuid(), 'Indeks Kedalaman Kemiskinan (P1)', 2022, 1.66),
(gen_random_uuid(), 'Indeks Kedalaman Kemiskinan (P1)', 2023, 2.29),
(gen_random_uuid(), 'Indeks Kedalaman Kemiskinan (P1)', 2024, 2.22),

-- Indeks Keparahan Kemiskinan (P2)
(gen_random_uuid(), 'Indeks Keparahan Kemiskinan (P2)', 2020, 0.55),
(gen_random_uuid(), 'Indeks Keparahan Kemiskinan (P2)', 2021, 0.47),
(gen_random_uuid(), 'Indeks Keparahan Kemiskinan (P2)', 2022, 0.31),
(gen_random_uuid(), 'Indeks Keparahan Kemiskinan (P2)', 2023, 0.56),
(gen_random_uuid(), 'Indeks Keparahan Kemiskinan (P2)', 2024, 0.61),

-- IPM
(gen_random_uuid(), 'IPM', 2020, 70.54),
(gen_random_uuid(), 'IPM', 2021, 71.04),
(gen_random_uuid(), 'IPM', 2022, 71.75),
(gen_random_uuid(), 'IPM', 2023, 72.47),
(gen_random_uuid(), 'IPM', 2024, 73.09),

-- Indeks Gini
(gen_random_uuid(), 'Indeks Gini', 2020, 0.34),
(gen_random_uuid(), 'Indeks Gini', 2021, 0.31),
(gen_random_uuid(), 'Indeks Gini', 2022, 0.30),
(gen_random_uuid(), 'Indeks Gini', 2023, 0.328),
(gen_random_uuid(), 'Indeks Gini', 2024, 0.289),

-- Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)
(gen_random_uuid(), 'Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)', 2020, 11418),
(gen_random_uuid(), 'Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)', 2021, 11459),
(gen_random_uuid(), 'Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)', 2022, 11563),
(gen_random_uuid(), 'Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)', 2023, 11897),
(gen_random_uuid(), 'Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)', 2024, 12414),

-- Umur Harapan Hidup (UHH)
(gen_random_uuid(), 'Umur Harapan Hidup (UHH)', 2020, 72.30),
(gen_random_uuid(), 'Umur Harapan Hidup (UHH)', 2021, 72.41),
(gen_random_uuid(), 'Umur Harapan Hidup (UHH)', 2022, 72.81),
(gen_random_uuid(), 'Umur Harapan Hidup (UHH)', 2023, 73.20),
(gen_random_uuid(), 'Umur Harapan Hidup (UHH)', 2024, 73.39),

-- Garis Kemiskinan (Rupiah)
(gen_random_uuid(), 'Garis Kemiskinan (Rupiah)', 2020, 342556),
(gen_random_uuid(), 'Garis Kemiskinan (Rupiah)', 2021, 358663),
(gen_random_uuid(), 'Garis Kemiskinan (Rupiah)', 2022, 382301),
(gen_random_uuid(), 'Garis Kemiskinan (Rupiah)', 2023, 413947),
(gen_random_uuid(), 'Garis Kemiskinan (Rupiah)', 2024, 445865);

-- +goose Down
DROP INDEX IF EXISTS idx_indikator_sosial_indikator_tahun;
DROP INDEX IF EXISTS idx_indikator_sosial_tahun;
DROP INDEX IF EXISTS idx_indikator_sosial_indikator;
DROP TABLE IF EXISTS indikator_sosial;