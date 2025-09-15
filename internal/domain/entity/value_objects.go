
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


type InstitutionType string

const (
    InstitutionDinas     InstitutionType = "DINAS"
    InstitutionDesa      InstitutionType = "DESA"
    InstitutionKecamatan InstitutionType = "KECAMATAN"
)


type AreaCategory string

const (
    AreaCagarBudaya            AreaCategory = "KAWASAN_CAGAR_BUDAYA"
    AreaHutan                  AreaCategory = "KAWASAN_HUTAN"
    AreaPariwisata             AreaCategory = "KAWASAN_PARIWISATA"
    AreaPerkebunan             AreaCategory = "KAWASAN_PERKEBUNAN"
    AreaPermukiman             AreaCategory = "KAWASAN_PERMUKIMAN"
    AreaPertahananKeamanan     AreaCategory = "KAWASAN_PERTAHANAN_KEAMANAN"
    AreaIndustri               AreaCategory = "KAWASAN_PERUNTUKAN_INDUSTRI"
    AreaPertambangan           AreaCategory = "KAWASAN_PERUNTUKAN_PERTAMBANGAN"
    AreaTanamanPangan          AreaCategory = "KAWASAN_TANAMAN_PANGAN"
    AreaTransportasi           AreaCategory = "KAWASAN_TRANSPORTASI"
    AreaLainnya                AreaCategory = "LAINNYA"
)


type SpatialViolationType string

const (
    ViolationSempadanSungai    SpatialViolationType = "BANGUNAN_SEMPADAN_SUNGAI"
    ViolationSempadanJalan     SpatialViolationType = "BANGUNAN_SEMPADAN_JALAN"
    ViolationAlihFungsiPertanian SpatialViolationType = "ALIH_FUNGSI_LAHAN_PERTANIAN"
    ViolationAlihFungsiRTH     SpatialViolationType = "ALIH_FUNGSI_RTH"
    ViolationTanpaIzin         SpatialViolationType = "PEMBANGUNAN_TANPA_IZIN"
    ViolationLainnya           SpatialViolationType = "LAINNYA"
)


type ViolationLevel string

const (
    ViolationRingan ViolationLevel = "RINGAN"
    ViolationSedang ViolationLevel = "SEDANG"
    ViolationBerat  ViolationLevel = "BERAT"
)


type EnvironmentalImpact string

const (
    ImpactKualitasRuang    EnvironmentalImpact = "MENURUN_KUALITAS_RUANG"
    ImpactBanjirLongsor    EnvironmentalImpact = "POTENSI_BANJIR_LONGSOR"
    ImpactGangguanAktivitas EnvironmentalImpact = "GANGGU_AKTIVITAS_WARGA"
)


type UrgencyLevel string

const (
    UrgencyMendesak UrgencyLevel = "MENDESAK"
    UrgencyBiasa    UrgencyLevel = "BIASA"
)


type SpatialReportStatus string

const (
    SpatialStatusPending    SpatialReportStatus = "PENDING"
    SpatialStatusReviewing  SpatialReportStatus = "REVIEWING"
    SpatialStatusProcessing SpatialReportStatus = "PROCESSING"
    SpatialStatusResolved   SpatialReportStatus = "RESOLVED"
    SpatialStatusRejected   SpatialReportStatus = "REJECTED"
)

type InstitutionUnitType string

const (
    InstitutionUnitDinas     InstitutionUnitType = "DINAS"
    InstitutionUnitDesa      InstitutionUnitType = "DESA"
    InstitutionUnitKecamatan InstitutionUnitType = "KECAMATAN"
)


type IrrigationType string

const (
    IrrigationSaluranSekunder IrrigationType = "SALURAN_SEKUNDER"
    IrrigationBendung         IrrigationType = "BENDUNG"
    IrrigationEmbungDam       IrrigationType = "EMBUNG_DAM"
    IrrigationPintuAir        IrrigationType = "PINTU_AIR"
)


type DamageType string

const (
    DamageRetakBocor           DamageType = "RETAK_BOCOR"
    DamageLongsorAmbrol        DamageType = "LONGSOR_AMBROL"
    DamageSedimentasiTinggi    DamageType = "SEDIMENTASI_TINGGI"
    DamageTersumbatSampah      DamageType = "TERSUMBAT_SAMPAH"
    DamageStrukturBetonRusak   DamageType = "STRUKTUR_BETON_RUSAK"
    DamagePintuAirMacet        DamageType = "PINTU_AIR_MACET"
    DamageTanggulJebol         DamageType = "TANGGUL_JEBOL"
    DamageLainnya              DamageType = "LAINNYA"
)


type DamageLevel string

const (
    DamageLevelRingan DamageLevel = "RINGAN"  
    DamageLevelSedang DamageLevel = "SEDANG"  
    DamageLevelBerat  DamageLevel = "BERAT"   
)


type UrgencyCategory string

const (
    UrgencyCategoryMendesak UrgencyCategory = "MENDESAK" 
    UrgencyCategoryRutin    UrgencyCategory = "RUTIN"
)


type WaterResourceStatus string

const (
    WaterResourceStatusPending     WaterResourceStatus = "PENDING"
    WaterResourceStatusVerified    WaterResourceStatus = "VERIFIED"
    WaterResourceStatusInProgress  WaterResourceStatus = "IN_PROGRESS"
    WaterResourceStatusCompleted   WaterResourceStatus = "COMPLETED"
    WaterResourceStatusPostponed   WaterResourceStatus = "POSTPONED"
    WaterResourceStatusRejected    WaterResourceStatus = "REJECTED"
)