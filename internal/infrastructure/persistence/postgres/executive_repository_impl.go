
package postgres

import (
    "context"
    "building-report-backend/internal/domain/entity"
    "building-report-backend/internal/domain/repository"
    
    "gorm.io/gorm"
)

type ExecutiveRepositoryImpl struct {
    db *gorm.DB
}

func NewExecutiveRepository(db *gorm.DB) repository.ExecutiveRepository {
    return &ExecutiveRepositoryImpl{db: db}
}

func (r *ExecutiveRepositoryImpl) FindByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorEkonomi, error) {
    var indikators []*entity.IndikatorEkonomi
    err := r.db.WithContext(ctx).
        Where("tahun = ?", tahun).
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorEkonomi, error) {
    var result entity.IndikatorEkonomi
    err := r.db.WithContext(ctx).
        Where("indikator = ? AND tahun = ?", indikator, tahun).
        First(&result).Error
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (r *ExecutiveRepositoryImpl) FindByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorEkonomi, error) {
    var indikators []*entity.IndikatorEkonomi
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorEkonomi, error) {
    var indikators []*entity.IndikatorEkonomi
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindDemografiByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorDemografi, error) {
    var indikators []*entity.IndikatorDemografi
    err := r.db.WithContext(ctx).
        Where("tahun = ?", tahun).
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindDemografiByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorDemografi, error) {
    var result entity.IndikatorDemografi
    err := r.db.WithContext(ctx).
        Where("indikator = ? AND tahun = ?", indikator, tahun).
        First(&result).Error
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (r *ExecutiveRepositoryImpl) FindDemografiByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorDemografi, error) {
    var indikators []*entity.IndikatorDemografi
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindDemografiAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorDemografi, error) {
    var indikators []*entity.IndikatorDemografi
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindSosialByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorSosial, error) {
    var indikators []*entity.IndikatorSosial
    err := r.db.WithContext(ctx).
        Where("tahun = ?", tahun).
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindSosialByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorSosial, error) {
    var result entity.IndikatorSosial
    err := r.db.WithContext(ctx).
        Where("indikator = ? AND tahun = ?", indikator, tahun).
        First(&result).Error
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (r *ExecutiveRepositoryImpl) FindSosialByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorSosial, error) {
    var indikators []*entity.IndikatorSosial
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindSosialAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorSosial, error) {
    var indikators []*entity.IndikatorSosial
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindKetenagakerjaanByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorKetenagakerjaan, error) {
    var indikators []*entity.IndikatorKetenagakerjaan
    err := r.db.WithContext(ctx).
        Where("tahun = ?", tahun).
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindKetenagakerjaanByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorKetenagakerjaan, error) {
    var result entity.IndikatorKetenagakerjaan
    err := r.db.WithContext(ctx).
        Where("indikator = ? AND tahun = ?", indikator, tahun).
        First(&result).Error
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (r *ExecutiveRepositoryImpl) FindKetenagakerjaanByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorKetenagakerjaan, error) {
    var indikators []*entity.IndikatorKetenagakerjaan
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindKetenagakerjaanAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorKetenagakerjaan, error) {
    var indikators []*entity.IndikatorKetenagakerjaan
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}


func (r *ExecutiveRepositoryImpl) FindPendidikanByTahun(ctx context.Context, tahun int) ([]*entity.IndikatorPendidikan, error) {
    var indikators []*entity.IndikatorPendidikan
    err := r.db.WithContext(ctx).
        Where("tahun = ?", tahun).
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindPendidikanByIndikatorAndTahun(ctx context.Context, indikator string, tahun int) (*entity.IndikatorPendidikan, error) {
    var result entity.IndikatorPendidikan
    err := r.db.WithContext(ctx).
        Where("indikator = ? AND tahun = ?", indikator, tahun).
        First(&result).Error
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (r *ExecutiveRepositoryImpl) FindPendidikanByIndikator(ctx context.Context, indikator string) ([]*entity.IndikatorPendidikan, error) {
    var indikators []*entity.IndikatorPendidikan
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}

func (r *ExecutiveRepositoryImpl) FindPendidikanAllTrend(ctx context.Context, indikator string) ([]*entity.IndikatorPendidikan, error) {
    var indikators []*entity.IndikatorPendidikan
    err := r.db.WithContext(ctx).
        Where("indikator = ?", indikator).
        Order("tahun ASC").
        Find(&indikators).Error
    if err != nil {
        return nil, err
    }
    return indikators, nil
}