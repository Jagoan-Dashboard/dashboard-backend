// internal/domain/entity/value_objects.go
package entity

type ReporterRole string

const (
    RolePerangkatDesa     ReporterRole = "PERANGKAT_DESA"
    RoleOPD               ReporterRole = "OPD_DINAS_TEKNIS"
    RoleKelompokMasyarakat ReporterRole = "KELOMPOK_MASYARAKAT"
    RoleMasyarakatUmum     ReporterRole = "MASYARAKAT_UMUM"
)

type BuildingType string

const (
    BuildingKantorPemerintah BuildingType = "KANTOR_PEMERINTAH"
    BuildingSekolah          BuildingType = "SEKOLAH"
    BuildingPuskesmas        BuildingType = "PUSKESMAS_POSYANDU"
    BuildingPasar            BuildingType = "PASAR"
    BuildingSaranaOlahraga   BuildingType = "SARANA_OLAHRAGA"
    BuildingFasilitasUmum    BuildingType = "FASILITAS_UMUM"
    BuildingLainnya          BuildingType = "LAINNYA"
)

type ReportStatusType string

const (
    StatusRehabilitasi   ReportStatusType = "REHABILITASI"
    StatusPembangunanBaru ReportStatusType = "PEMBANGUNAN_BARU"
    StatusLainnya        ReportStatusType = "LAINNYA"
)

type FundingSource string

const (
    FundingAPBDKab     FundingSource = "APBD_KABUPATEN"
    FundingAPBDProv    FundingSource = "APBD_PROVINSI"
    FundingAPBN        FundingSource = "APBN"
    FundingDanaDesa    FundingSource = "DANA_DESA"
    FundingSwadaya     FundingSource = "SWADAYA_MASYARAKAT"
    FundingLainnya     FundingSource = "LAINNYA"
)

type WorkType string

const (
    WorkPerbaikanAtap      WorkType = "PERBAIKAN_ATAP"
    WorkPerbaikanDinding   WorkType = "PERBAIKAN_DINDING"
    WorkPerbaikanLantai    WorkType = "PERBAIKAN_LANTAI"
    WorkPerbaikanPintu     WorkType = "PERBAIKAN_PINTU_JENDELA"
    WorkPerbaikanSanitasi  WorkType = "PERBAIKAN_SANITASI"
    WorkPerbaikanListrik   WorkType = "PERBAIKAN_LISTRIK_AIR"
    WorkLainnya           WorkType = "LAINNYA"
)

type ConditionAfterRehab string

const (
    ConditionBaik              ConditionAfterRehab = "BAIK_SIAP_PAKAI"
    ConditionButuhPerbaikan    ConditionAfterRehab = "BUTUH_PERBAIKAN_TAMBAHAN"
    ConditionLainnya           ConditionAfterRehab = "LAINNYA"
)