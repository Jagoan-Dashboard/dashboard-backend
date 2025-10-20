package usecase

import (
	"context"
	"errors"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"

	"gorm.io/gorm"
)

var (
    ErrDataNotFound = errors.New("data tidak ditemukan untuk tahun yang diminta")
)

type ExecutiveUseCase struct {
    executiveRepo repository.ExecutiveRepository
}

func NewExecutiveUseCase(executiveRepo repository.ExecutiveRepository) *ExecutiveUseCase {
    return &ExecutiveUseCase{
        executiveRepo: executiveRepo,
    }
}

func (uc *ExecutiveUseCase) GetEkonomiOverview(ctx context.Context, tahun int) (*dto.EkonomiOverviewResponse, error) {
    
    dataCurrentYear, err := uc.executiveRepo.FindByTahun(ctx, tahun)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &dto.EkonomiOverviewResponse{
                Tahun:                tahun,
                TrendLajuPertumbuhan: []dto.TrendData{},
                TrendInflasi:         []dto.TrendData{},
            }, nil
        }
        return nil, err
    }

    if len(dataCurrentYear) == 0 {
        return &dto.EkonomiOverviewResponse{
            Tahun:                tahun,
            TrendLajuPertumbuhan: []dto.TrendData{},
            TrendInflasi:         []dto.TrendData{},
        }, nil
    }

    
    dataPreviousYear, _ := uc.executiveRepo.FindByTahun(ctx, tahun-1)

    
    currentYearMap := make(map[string]float64)
    for _, data := range dataCurrentYear {
        currentYearMap[data.Indikator] = data.Nilai
    }

    previousYearMap := make(map[string]float64)
    for _, data := range dataPreviousYear {
        previousYearMap[data.Indikator] = data.Nilai
    }

    
    trendLPE, _ := uc.executiveRepo.FindAllTrend(ctx, "Laju Pertumbuhan Ekonomi")
    trendInflasi, _ := uc.executiveRepo.FindAllTrend(ctx, "Inflasi")

    
    response := &dto.EkonomiOverviewResponse{
        Tahun:                  tahun,
        LajuPertumbuhanEkonomi: currentYearMap["Laju Pertumbuhan Ekonomi"],
        PertanianPDRB:          currentYearMap["% Pertanian (PDRB)"],
        PengolahanPDRB:         currentYearMap["% Pengolahan (PDRB)"],
        ICOR:                   currentYearMap["ICOR"],
        ILOR:                   currentYearMap["ILOR"],
        Inflasi:                currentYearMap["Inflasi"],
        TrendLajuPertumbuhan:   uc.convertToTrendData(trendLPE),
        TrendInflasi:           uc.convertToTrendData(trendInflasi),
    }

    
    response.PerubahanLPE = uc.calculatePercentageChange(
        previousYearMap["Laju Pertumbuhan Ekonomi"],
        currentYearMap["Laju Pertumbuhan Ekonomi"],
    )
    response.PerubahanPertanian = uc.calculatePercentageChange(
        previousYearMap["% Pertanian (PDRB)"],
        currentYearMap["% Pertanian (PDRB)"],
    )
    response.PerubahanPengolahan = uc.calculatePercentageChange(
        previousYearMap["% Pengolahan (PDRB)"],
        currentYearMap["% Pengolahan (PDRB)"],
    )
    response.PerubahanICOR = uc.calculatePercentageChange(
        previousYearMap["ICOR"],
        currentYearMap["ICOR"],
    )
    response.PerubahanILOR = uc.calculatePercentageChange(
        previousYearMap["ILOR"],
        currentYearMap["ILOR"],
    )
    response.PerubahanInflasi = uc.calculatePercentageChange(
        previousYearMap["Inflasi"],
        currentYearMap["Inflasi"],
    )

    return response, nil
}

func (uc *ExecutiveUseCase) calculatePercentageChange(oldValue, newValue float64) *float64 {
    if oldValue == 0 {
        return nil
    }
    
    change := ((newValue - oldValue) / oldValue) * 100
    return &change
}

func (uc *ExecutiveUseCase) convertToTrendData(data []*entity.IndikatorEkonomi) []dto.TrendData {
    result := make([]dto.TrendData, 0, len(data))
    for _, item := range data {
        result = append(result, dto.TrendData{
            Tahun: item.Tahun,
            Nilai: item.Nilai,
        })
    }
    return result
}

func (uc *ExecutiveUseCase) GetPopulationOverview(ctx context.Context, tahun int) (*dto.PopulationOverviewResponse, error) {
    dataCurrentYear, err := uc.executiveRepo.FindDemografiByTahun(ctx, tahun)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &dto.PopulationOverviewResponse{
                Tahun:                    tahun,
                TrendKepadatanPenduduk:   []dto.TrendData{},
                TrendRasioKetergantungan: []dto.TrendData{},
            }, nil
        }
        return nil, err
    }

    if len(dataCurrentYear) == 0 {
        return &dto.PopulationOverviewResponse{
            Tahun:                    tahun,
            TrendKepadatanPenduduk:   []dto.TrendData{},
            TrendRasioKetergantungan: []dto.TrendData{},
        }, nil
    }

    dataPreviousYear, _ := uc.executiveRepo.FindDemografiByTahun(ctx, tahun-1)

    currentYearMap := make(map[string]float64)
    for _, data := range dataCurrentYear {
        currentYearMap[data.Indikator] = data.Nilai
    }

    previousYearMap := make(map[string]float64)
    for _, data := range dataPreviousYear {
        previousYearMap[data.Indikator] = data.Nilai
    }

    trendKepadatan, _ := uc.executiveRepo.FindDemografiAllTrend(ctx, "Kepadatan Penduduk")
    trendRasio, _ := uc.executiveRepo.FindDemografiAllTrend(ctx, "Rasio Ketergantungan")

    response := &dto.PopulationOverviewResponse{
        Tahun:                    tahun,
        KepadatanPenduduk:        currentYearMap["Kepadatan Penduduk"],
        RasioKetergantungan:      currentYearMap["Rasio Ketergantungan"],
        PendudukProduktif:        currentYearMap["Jumlah penduduk produktif (Usia 15–64 tahun)"],
        PendudukNonProduktif:     currentYearMap["Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)"],
        TrendKepadatanPenduduk:   uc.convertToTrendDataDemografi(trendKepadatan),
        TrendRasioKetergantungan: uc.convertToTrendDataDemografi(trendRasio),
    }

    response.PerubahanKepadatan = uc.calculatePercentageChange(
        previousYearMap["Kepadatan Penduduk"],
        currentYearMap["Kepadatan Penduduk"],
    )
    response.PerubahanRasio = uc.calculatePercentageChange(
        previousYearMap["Rasio Ketergantungan"],
        currentYearMap["Rasio Ketergantungan"],
    )
    response.PerubahanProduktif = uc.calculatePercentageChange(
        previousYearMap["Jumlah penduduk produktif (Usia 15–64 tahun)"],
        currentYearMap["Jumlah penduduk produktif (Usia 15–64 tahun)"],
    )
    response.PerubahanNonProduktif = uc.calculatePercentageChange(
        previousYearMap["Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)"],
        currentYearMap["Jumlah penduduk non produktif (Usia <15 Tahun dan Usia 65 Tahun ke atas)"],
    )

    return response, nil
}

func (uc *ExecutiveUseCase) convertToTrendDataEkonomi(data []*entity.IndikatorEkonomi) []dto.TrendData {
    result := make([]dto.TrendData, 0, len(data))
    for _, item := range data {
        result = append(result, dto.TrendData{
            Tahun: item.Tahun,
            Nilai: item.Nilai,
        })
    }
    return result
}

func (uc *ExecutiveUseCase) convertToTrendDataDemografi(data []*entity.IndikatorDemografi) []dto.TrendData {
    result := make([]dto.TrendData, 0, len(data))
    for _, item := range data {
        result = append(result, dto.TrendData{
            Tahun: item.Tahun,
            Nilai: item.Nilai,
        })
    }
    return result
}

func (uc *ExecutiveUseCase) GetPovertyOverview(ctx context.Context, tahun int) (*dto.PovertyOverviewResponse, error) {
    dataCurrentYear, err := uc.executiveRepo.FindSosialByTahun(ctx, tahun)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &dto.PovertyOverviewResponse{
                Tahun:                          tahun,
                TrendIndeksKedalamanKemiskinan: []dto.TrendData{},
                TrendIndeksKeparahanKemiskinan: []dto.TrendData{},
            }, nil
        }
        return nil, err
    }

    if len(dataCurrentYear) == 0 {
        return &dto.PovertyOverviewResponse{
            Tahun:                          tahun,
            TrendIndeksKedalamanKemiskinan: []dto.TrendData{},
            TrendIndeksKeparahanKemiskinan: []dto.TrendData{},
        }, nil
    }

    dataPreviousYear, _ := uc.executiveRepo.FindSosialByTahun(ctx, tahun-1)

    currentYearMap := make(map[string]float64)
    for _, data := range dataCurrentYear {
        currentYearMap[data.Indikator] = data.Nilai
    }

    previousYearMap := make(map[string]float64)
    for _, data := range dataPreviousYear {
        previousYearMap[data.Indikator] = data.Nilai
    }

    trendKedalaman, _ := uc.executiveRepo.FindSosialAllTrend(ctx, "Indeks Kedalaman Kemiskinan (P1)")
    trendKeparahan, _ := uc.executiveRepo.FindSosialAllTrend(ctx, "Indeks Keparahan Kemiskinan (P2)")

    response := &dto.PovertyOverviewResponse{
        Tahun:                          tahun,
        AngkaKemiskinan:                currentYearMap["Angka Kemiskinan (P0)"],
        IndeksKedalamanKemiskinan:      currentYearMap["Indeks Kedalaman Kemiskinan (P1)"],
        IndeksKeparahanKemiskinan:      currentYearMap["Indeks Keparahan Kemiskinan (P2)"],
        IPM:                            currentYearMap["IPM"],
        IndeksGini:                     currentYearMap["Indeks Gini"],
        PengeluaranPerKapita:           currentYearMap["Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)"],
        UmurHarapanHidup:               currentYearMap["Umur Harapan Hidup (UHH)"],
        GarisKemiskinan:                currentYearMap["Garis Kemiskinan (Rupiah)"],
        TrendIndeksKedalamanKemiskinan: uc.convertToTrendDataSosial(trendKedalaman),
        TrendIndeksKeparahanKemiskinan: uc.convertToTrendDataSosial(trendKeparahan),
    }

    response.PerubahanAngkaKemiskinan = uc.calculatePercentageChange(
        previousYearMap["Angka Kemiskinan (P0)"],
        currentYearMap["Angka Kemiskinan (P0)"],
    )
    response.PerubahanIndeksKedalaman = uc.calculatePercentageChange(
        previousYearMap["Indeks Kedalaman Kemiskinan (P1)"],
        currentYearMap["Indeks Kedalaman Kemiskinan (P1)"],
    )
    response.PerubahanIndeksKeparahan = uc.calculatePercentageChange(
        previousYearMap["Indeks Keparahan Kemiskinan (P2)"],
        currentYearMap["Indeks Keparahan Kemiskinan (P2)"],
    )
    response.PerubahanIPM = uc.calculatePercentageChange(
        previousYearMap["IPM"],
        currentYearMap["IPM"],
    )
    response.PerubahanIndeksGini = uc.calculatePercentageChange(
        previousYearMap["Indeks Gini"],
        currentYearMap["Indeks Gini"],
    )
    response.PerubahanPengeluaran = uc.calculatePercentageChange(
        previousYearMap["Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)"],
        currentYearMap["Pengeluaran Per Kapita Riil Disesuaikan (Ribu Rupiah)"],
    )
    response.PerubahanUHH = uc.calculatePercentageChange(
        previousYearMap["Umur Harapan Hidup (UHH)"],
        currentYearMap["Umur Harapan Hidup (UHH)"],
    )
    response.PerubahanGarisKemiskinan = uc.calculatePercentageChange(
        previousYearMap["Garis Kemiskinan (Rupiah)"],
        currentYearMap["Garis Kemiskinan (Rupiah)"],
    )

    return response, nil
}

func (uc *ExecutiveUseCase) convertToTrendDataSosial(data []*entity.IndikatorSosial) []dto.TrendData {
    result := make([]dto.TrendData, 0, len(data))
    for _, item := range data {
        result = append(result, dto.TrendData{
            Tahun: item.Tahun,
            Nilai: item.Nilai,
        })
    }
    return result
}

func (uc *ExecutiveUseCase) GetEmploymentOverview(ctx context.Context, tahun int) (*dto.EmploymentOverviewResponse, error) {
    dataCurrentYear, err := uc.executiveRepo.FindKetenagakerjaanByTahun(ctx, tahun)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &dto.EmploymentOverviewResponse{
                Tahun:     tahun,
                TrendTPT:  []dto.TrendData{},
                TrendTPAK: []dto.TrendData{},
            }, nil
        }
        return nil, err
    }

    if len(dataCurrentYear) == 0 {
        return &dto.EmploymentOverviewResponse{
            Tahun:     tahun,
            TrendTPT:  []dto.TrendData{},
            TrendTPAK: []dto.TrendData{},
        }, nil
    }

    dataPreviousYear, _ := uc.executiveRepo.FindKetenagakerjaanByTahun(ctx, tahun-1)

    currentYearMap := make(map[string]float64)
    for _, data := range dataCurrentYear {
        currentYearMap[data.Indikator] = data.Nilai
    }

    previousYearMap := make(map[string]float64)
    for _, data := range dataPreviousYear {
        previousYearMap[data.Indikator] = data.Nilai
    }

    trendTPT, _ := uc.executiveRepo.FindKetenagakerjaanAllTrend(ctx, "TPT")
    trendTPAK, _ := uc.executiveRepo.FindKetenagakerjaanAllTrend(ctx, "TPAK")

    response := &dto.EmploymentOverviewResponse{
        Tahun:         tahun,
        TPT:           currentYearMap["TPT"],
        TPAK:          currentYearMap["TPAK"],
        TPAKPerempuan: currentYearMap["TPAK Perempuan"],
        UpahMinimum:   currentYearMap["Upah Minimum Kabupaten (Rupiah)"],
        TrendTPT:      uc.convertToTrendDataKetenagakerjaan(trendTPT),
        TrendTPAK:     uc.convertToTrendDataKetenagakerjaan(trendTPAK),
    }

    response.PerubahanTPT = uc.calculatePercentageChange(
        previousYearMap["TPT"],
        currentYearMap["TPT"],
    )
    response.PerubahanTPAK = uc.calculatePercentageChange(
        previousYearMap["TPAK"],
        currentYearMap["TPAK"],
    )
    response.PerubahanTPAKPerempuan = uc.calculatePercentageChange(
        previousYearMap["TPAK Perempuan"],
        currentYearMap["TPAK Perempuan"],
    )
    response.PerubahanUpahMinimum = uc.calculatePercentageChange(
        previousYearMap["Upah Minimum Kabupaten (Rupiah)"],
        currentYearMap["Upah Minimum Kabupaten (Rupiah)"],
    )

    return response, nil
}

func (uc *ExecutiveUseCase) convertToTrendDataKetenagakerjaan(data []*entity.IndikatorKetenagakerjaan) []dto.TrendData {
    result := make([]dto.TrendData, 0, len(data))
    for _, item := range data {
        result = append(result, dto.TrendData{
            Tahun: item.Tahun,
            Nilai: item.Nilai,
        })
    }
    return result
}

func (uc *ExecutiveUseCase) GetEducationOverview(ctx context.Context, tahun int) (*dto.EducationOverviewResponse, error) {
    dataCurrentYear, err := uc.executiveRepo.FindPendidikanByTahun(ctx, tahun)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &dto.EducationOverviewResponse{
                Tahun:                    tahun,
                TrendRataRataLamaSekolah: []dto.TrendData{},
                TrendHarapanLamaSekolah:  []dto.TrendData{},
            }, nil
        }
        return nil, err
    }

    if len(dataCurrentYear) == 0 {
        return &dto.EducationOverviewResponse{
            Tahun:                    tahun,
            TrendRataRataLamaSekolah: []dto.TrendData{},
            TrendHarapanLamaSekolah:  []dto.TrendData{},
        }, nil
    }

    dataPreviousYear, _ := uc.executiveRepo.FindPendidikanByTahun(ctx, tahun-1)

    currentYearMap := make(map[string]float64)
    for _, data := range dataCurrentYear {
        currentYearMap[data.Indikator] = data.Nilai
    }

    previousYearMap := make(map[string]float64)
    for _, data := range dataPreviousYear {
        previousYearMap[data.Indikator] = data.Nilai
    }

    trendRataRata, _ := uc.executiveRepo.FindPendidikanAllTrend(ctx, "Rata-rata Lama Sekolah")
    trendHarapan, _ := uc.executiveRepo.FindPendidikanAllTrend(ctx, "Harapan Lama Sekolah")

    response := &dto.EducationOverviewResponse{
        Tahun:                    tahun,
        RataRataLamaSekolah:      currentYearMap["Rata-rata Lama Sekolah"],
        HarapanLamaSekolah:       currentYearMap["Harapan Lama Sekolah"],
        ProporsiPendidikanTinggi: currentYearMap["Proporsi dengan Pendidikan Tinggi"],
        TrendRataRataLamaSekolah: uc.convertToTrendDataPendidikan(trendRataRata),
        TrendHarapanLamaSekolah:  uc.convertToTrendDataPendidikan(trendHarapan),
    }

    response.PerubahanRataRataLamaSekolah = uc.calculatePercentageChange(
        previousYearMap["Rata-rata Lama Sekolah"],
        currentYearMap["Rata-rata Lama Sekolah"],
    )
    response.PerubahanHarapanLamaSekolah = uc.calculatePercentageChange(
        previousYearMap["Harapan Lama Sekolah"],
        currentYearMap["Harapan Lama Sekolah"],
    )
    response.PerubahanProporsiPendidikan = uc.calculatePercentageChange(
        previousYearMap["Proporsi dengan Pendidikan Tinggi"],
        currentYearMap["Proporsi dengan Pendidikan Tinggi"],
    )

    return response, nil
}

func (uc *ExecutiveUseCase) convertToTrendDataPendidikan(data []*entity.IndikatorPendidikan) []dto.TrendData {
    result := make([]dto.TrendData, 0, len(data))
    for _, item := range data {
        result = append(result, dto.TrendData{
            Tahun: item.Tahun,
            Nilai: item.Nilai,
        })
    }
    return result
}