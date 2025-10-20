
package repository

import (
    "context"
    "building-report-backend/internal/domain/entity"
)

type ExecutiveRepository interface {
    // Ekonomi methods
    FindByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorEkonomi, error)
    FindByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorEkonomi, error)
    FindByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorEkonomi, error)
    FindAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorEkonomi, error)

     // Demografi methods
    FindDemografiByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorDemografi, error)
    FindDemografiByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorDemografi, error)
    FindDemografiByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorDemografi, error)
    FindDemografiAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorDemografi, error)

    // Sosial methods
    FindSosialByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorSosial, error)
    FindSosialByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorSosial, error)
    FindSosialByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorSosial, error)
    FindSosialAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorSosial, error)

    // Ketenagakerjaan methods
    FindKetenagakerjaanByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorKetenagakerjaan, error)
    FindKetenagakerjaanByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorKetenagakerjaan, error)
    FindKetenagakerjaanByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorKetenagakerjaan, error)
    FindKetenagakerjaanAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorKetenagakerjaan, error)

    // Pendidikan methods
    FindPendidikanByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorPendidikan, error)
    FindPendidikanByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorPendidikan, error)
    FindPendidikanByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorPendidikan, error)
    FindPendidikanAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorPendidikan, error)
}