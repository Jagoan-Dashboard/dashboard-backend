
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

type RoadType string

const (
    RoadTypeJalanNasional   RoadType = "JALAN_NASIONAL"
    RoadTypeJalanProvinsi   RoadType = "JALAN_PROVINSI"
    RoadTypeJalanKabupaten  RoadType = "JALAN_KABUPATEN"
    RoadTypeJalanKota       RoadType = "JALAN_KOTA"
    RoadTypeJalanDesa       RoadType = "JALAN_DESA"
    RoadTypeJalanLingkungan RoadType = "JALAN_LINGKUNGAN"
)


type RoadClass string

const (
    RoadClassArteri      RoadClass = "ARTERI"      
    RoadClassKolektor    RoadClass = "KOLEKTOR"    
    RoadClassLokal       RoadClass = "LOKAL"       
    RoadClassLingkungan  RoadClass = "LINGKUNGAN"  
)


type RoadDamageType string

const (
    RoadDamageRetak         RoadDamageType = "RETAK_MEMANJANG"      
    RoadDamageRetakMelintang RoadDamageType = "RETAK_MELINTANG"     
    RoadDamageRetakBlok     RoadDamageType = "RETAK_BLOK"          
    RoadDamageRetakBuaya    RoadDamageType = "RETAK_KULIT_BUAYA"   
    RoadDamageLubang        RoadDamageType = "LUBANG"              
    RoadDamageAmblas        RoadDamageType = "AMBLAS"              
    RoadDamageGelombang     RoadDamageType = "GELOMBANG"           
    RoadDamageTepiJalan     RoadDamageType = "KERUSAKAN_TEPI"      
    RoadDamageDrainase      RoadDamageType = "KERUSAKAN_DRAINASE"  
    RoadDamageJembatan      RoadDamageType = "KERUSAKAN_JEMBATAN"  
    RoadDamagePerlengkapan  RoadDamageType = "KERUSAKAN_PERLENGKAPAN" 
    RoadDamageLainnya       RoadDamageType = "LAINNYA"
)


type RoadDamageLevel string

const (
    RoadDamageLevelMinor    RoadDamageLevel = "RINGAN"   
    RoadDamageLevelModerate RoadDamageLevel = "SEDANG"   
    RoadDamageLevelSevere   RoadDamageLevel = "BERAT"    
)


type TrafficImpact string

const (
    TrafficImpactMinimal         TrafficImpact = "MINIMAL"           
    TrafficImpactReduced         TrafficImpact = "TERGANGGU"         
    TrafficImpactSeverelyReduced TrafficImpact = "SANGAT_TERGANGGU"  
    TrafficImpactBlocked         TrafficImpact = "TERPUTUS"          
)


type RoadUrgencyLevel string

const (
    RoadUrgencyLow       RoadUrgencyLevel = "RENDAH"    
    RoadUrgencyMedium    RoadUrgencyLevel = "SEDANG"    
    RoadUrgencyHigh      RoadUrgencyLevel = "TINGGI"    
    RoadUrgencyEmergency RoadUrgencyLevel = "DARURAT"   
)


type BinaMargaStatus string

const (
    BinaMargaStatusPending     BinaMargaStatus = "PENDING"      
    BinaMargaStatusVerified    BinaMargaStatus = "VERIFIED"     
    BinaMargaStatusPlanned     BinaMargaStatus = "PLANNED"      
    BinaMargaStatusInProgress  BinaMargaStatus = "IN_PROGRESS"  
    BinaMargaStatusCompleted   BinaMargaStatus = "COMPLETED"    
    BinaMargaStatusPostponed   BinaMargaStatus = "POSTPONED"    
    BinaMargaStatusRejected    BinaMargaStatus = "REJECTED"     
)

type FarmerGroupType string

const (
    FarmerGroupPoktan  FarmerGroupType = "POKTAN"   
    FarmerGroupGapoktan FarmerGroupType = "GAPOKTAN" 
)


type FoodCommodity string

const (
    FoodCommodityPadiSawah   FoodCommodity = "PADI_SAWAH"
    FoodCommodityPadiLadang  FoodCommodity = "PADI_LADANG"
    FoodCommodityJagung      FoodCommodity = "JAGUNG"
    FoodCommodityKedelai     FoodCommodity = "KEDELAI"
    FoodCommodityKacangTanah FoodCommodity = "KACANG_TANAH"
    FoodCommodityUbiKayu     FoodCommodity = "UBI_KAYU"
    FoodCommodityUbiJalar    FoodCommodity = "UBI_JALAR"
    FoodCommodityLainnya     FoodCommodity = "LAINNYA"
)


type HorticultureCommodity string

const (
    HortiCommoditySayuran      HorticultureCommodity = "SAYURAN"
    HortiCommodityBuah         HorticultureCommodity = "BUAH"
    HortiCommodityFlorikultura HorticultureCommodity = "FLORIKULTURA"
)


type PlantationCommodity string

const (
    PlantationCommodityKopi       PlantationCommodity = "KOPI"
    PlantationCommodityKakao      PlantationCommodity = "KAKAO"
    PlantationCommodityKelapa     PlantationCommodity = "KELAPA"
    PlantationCommodityKelapaSawit PlantationCommodity = "KELAPA_SAWIT"
    PlantationCommodityCengkeh    PlantationCommodity = "CENGKEH"
    PlantationCommodityTebu       PlantationCommodity = "TEBU"
    PlantationCommodityKaret      PlantationCommodity = "KARET"
    PlantationCommodityTembakau   PlantationCommodity = "TEMBAKAU"
    PlantationCommodityLainnya    PlantationCommodity = "LAINNYA"
)


type LandStatus string

const (
    LandStatusOwned  LandStatus = "MILIK_SENDIRI"
    LandStatusRented LandStatus = "SEWA"
    LandStatusFree   LandStatus = "PINJAM_BEBAS_SEWA"
    LandStatusOther  LandStatus = "LAINNYA"
)


type GrowthPhase string

const (
    GrowthPhaseBelumTanam     GrowthPhase = "BELUM_TANAM"
    GrowthPhaseBera           GrowthPhase = "BERA"
    GrowthPhaseVegetatifAwal  GrowthPhase = "VEGETATIF_AWAL"
    GrowthPhaseVegetatifAkhir GrowthPhase = "VEGETATIF_AKHIR"
    GrowthPhaseGeneratif1     GrowthPhase = "GENERATIF_1"
    GrowthPhaseGeneratif2     GrowthPhase = "GENERATIF_2"
    GrowthPhaseGeneratif3     GrowthPhase = "GENERATIF_3"
    GrowthPhasePanenMuda      GrowthPhase = "PANEN_MUDA"
    GrowthPhasePanenPenuh     GrowthPhase = "PANEN_PENUH"
    GrowthPhaseLainnya        GrowthPhase = "LAINNYA"
)


type HortiGrowthPhase string

const (
    HortiGrowthPhasePersemaian   HortiGrowthPhase = "PERSEMAIAN"
    HortiGrowthPhaseVegetatif    HortiGrowthPhase = "VEGETATIF"
    HortiGrowthPhasePembungaan   HortiGrowthPhase = "PEMBUNGAAN"
    HortiGrowthPhasePembuahan    HortiGrowthPhase = "PEMBUAHAN"
    HortiGrowthPhasePanen        HortiGrowthPhase = "PANEN"
    HortiGrowthPhaseLainnya      HortiGrowthPhase = "LAINNYA"
)


type PlantationGrowthPhase string

const (
    PlantationGrowthPhaseBibit             PlantationGrowthPhase = "BIBIT_PERSEMAIAN"
    PlantationGrowthPhaseTanamanMuda       PlantationGrowthPhase = "TANAMAN_MUDA_TBM"
    PlantationGrowthPhaseTanamanMenghasilkan PlantationGrowthPhase = "TANAMAN_MENGHASILKAN_TM"
    PlantationGrowthPhasePanen             PlantationGrowthPhase = "PANEN"
    PlantationGrowthPhaseLainnya           PlantationGrowthPhase = "LAINNYA"
)


type DelayReason string

const (
    DelayReasonTidakAda     DelayReason = "TIDAK_ADA"
    DelayReasonHujanTerus   DelayReason = "ADA_HUJAN_TERUS"
    DelayReasonKekeringan   DelayReason = "ADA_KEKERINGAN"
    DelayReasonBibitTerlambat DelayReason = "ADA_BIBIT_TERLAMBAT"
    DelayReasonBanjir       DelayReason = "ADA_BANJIR"
    DelayReasonLainnya      DelayReason = "LAINNYA"
)


type TechnologyMethod string

const (
    TechnologyMethodTidakAda      TechnologyMethod = "TIDAK_ADA"
    TechnologyMethodJajarLegowo   TechnologyMethod = "JAJAR_LEGOWO"
    TechnologyMethodDroneSemprot  TechnologyMethod = "DRONE_SEMPROT"
    TechnologyMethodPupukOrganik  TechnologyMethod = "PUPUK_ORGANIK_HAYATI"
    TechnologyMethodIrigasiPompa  TechnologyMethod = "IRIGASI_POMPA"
    TechnologyMethodLainnya       TechnologyMethod = "LAINNYA"
)


type HortiTechnology string

const (
    HortiTechnologyTidakAda       HortiTechnology = "TIDAK_ADA"
    HortiTechnologyGreenhouse     HortiTechnology = "GREENHOUSE_SCREEN_HOUSE"
    HortiTechnologyMulsaPlastik   HortiTechnology = "MULSA_PLASTIK"
    HortiTechnologyIrigasiTetes   HortiTechnology = "IRIGASI_TETES_SPRINKLER"
    HortiTechnologyDroneSemprot   HortiTechnology = "DRONE_SEMPROT"
    HortiTechnologySensorIoT      HortiTechnology = "SENSOR_IOT_KELEMBAPAN"
    HortiTechnologyPupukOrganik   HortiTechnology = "PUPUK_ORGANIK_HAYATI"
    HortiTechnologyLainnya        HortiTechnology = "LAINNYA"
)


type PlantationTechnology string

const (
    PlantationTechnologyTidakAda        PlantationTechnology = "TIDAK_ADA"
    PlantationTechnologyPeremajaan      PlantationTechnology = "PEREMAJAAN_TANAMAN_REPLANTING"
    PlantationTechnologyPupukOrganik    PlantationTechnology = "PUPUK_ORGANIK_HAYATI"
    PlantationTechnologyIrigasiTetes    PlantationTechnology = "IRIGASI_TETES_SPRINKLER"
    PlantationTechnologyDroneMonitoring PlantationTechnology = "DRONE_MONITORING"
    PlantationTechnologyAgroforestry    PlantationTechnology = "SISTEM_AGROFORESTRY"
    PlantationTechnologyLainnya         PlantationTechnology = "LAINNYA"
)


type PostHarvestProblem string

const (
    PostHarvestProblemTidakAda         PostHarvestProblem = "TIDAK_ADA"
    PostHarvestProblemSusutTinggi      PostHarvestProblem = "SUSUT_TINGGI_HASIL_CEPAT_BUSUK"
    PostHarvestProblemKeterbatasanGudang PostHarvestProblem = "KETERBATASAN_GUDANG_PENDINGIN_COLD_STORAGE"
    PostHarvestProblemKesulitanKemasan PostHarvestProblem = "KESULITAN_KEMASAN_GRADING"
    PostHarvestProblemKesulitanAkses   PostHarvestProblem = "KESULITAN_AKSES_PASAR"
    PostHarvestProblemLainnya          PostHarvestProblem = "LAINNYA"
)


type ProductionProblem string

const (
    ProductionProblemRendahnyaProduktivitas ProductionProblem = "RENDAHNYA_PRODUKTIVITAS"
    ProductionProblemHargaJualFluktuatif    ProductionProblem = "HARGA_JUAL_FLUKTUATIF"
    ProductionProblemSeranganHama           ProductionProblem = "SERANGAN_HAMA_PENYAKIT"
    ProductionProblemPerluReplanting        ProductionProblem = "PERLU_REPLANTING"
    ProductionProblemLainnya                ProductionProblem = "LAINNYA"
)


type PestDiseaseType string

const (
    
    PestDiseaseUlatGrayak  PestDiseaseType = "ULAT_GRAYAK"
    PestDiseaseWerengCoklat PestDiseaseType = "WERENG_COKLAT"
    PestDiseaseTikus       PestDiseaseType = "TIKUS"
    PestDiseaseBusukDaun   PestDiseaseType = "BUSUK_DAUN"
    
    
    PestDiseaseTrips        PestDiseaseType = "TRIPS"
    PestDiseaseLalatBuah    PestDiseaseType = "LALAT_BUAH"
    PestDiseaseAntraknosa   PestDiseaseType = "ANTRAKNOSA"
    PestDiseaseLayuFusarium PestDiseaseType = "LAYU_FUSARIUM"
    
    
    PestDiseasePBKo        PestDiseaseType = "PBKO_PENGGEREK_BUAH_KAKAO"
    PestDiseaseKaratDaun   PestDiseaseType = "KARAT_DAUN"
    PestDiseaseBusukBatang PestDiseaseType = "BUSUK_BATANG"
    PestDiseaseHamaTikus   PestDiseaseType = "HAMA_TIKUS"
    
    PestDiseaseLainnya     PestDiseaseType = "LAINNYA"
)


type AffectedAreaLevel string

const (
    AffectedAreaLevelKurang10   AffectedAreaLevel = "KURANG_10_PERSEN"
    AffectedAreaLevel10Sampai25 AffectedAreaLevel = "10_SAMPAI_25_PERSEN"
    AffectedAreaLevel25Sampai50 AffectedAreaLevel = "25_SAMPAI_50_PERSEN"
    AffectedAreaLevelLebih50    AffectedAreaLevel = "LEBIH_50_PERSEN"
    AffectedAreaLevelLainnya    AffectedAreaLevel = "LAINNYA"
)


type ControlAction string

const (
    ControlActionBelumDitangani     ControlAction = "BELUM_DITANGANI"
    ControlActionSemprotInsektisida ControlAction = "SEMPROT_INSEKTISIDA_KIMIA"
    ControlActionSemprotBiopestisida ControlAction = "SEMPROT_BIOPESTISIDA"
    ControlActionPasangPerangkap    ControlAction = "PASANG_PERANGKAP_SETRUM"
    ControlActionSanitasiKebun      ControlAction = "SANITASI_KEBUN_ROTASI_TANAMAN"
    ControlActionLainnya            ControlAction = "LAINNYA"
)


type WeatherCondition string

const (
    WeatherConditionHujan        WeatherCondition = "HUJAN"
    WeatherConditionCerah        WeatherCondition = "CERAH"
    WeatherConditionAnginKencang WeatherCondition = "ANGIN_KENCANG"
    WeatherConditionKekeringan   WeatherCondition = "KEKERINGAN"
    WeatherConditionLainnya      WeatherCondition = "LAINNYA"
)


type WeatherImpact string

const (
    WeatherImpactTidakAda      WeatherImpact = "TIDAK_ADA"
    WeatherImpactTanamanRebah  WeatherImpact = "TANAMAN_REBAH"
    WeatherImpactDaunMenguning WeatherImpact = "DAUN_MENGUNING"
    WeatherImpactBuahRontok    WeatherImpact = "BUAH_RONTOK"
    WeatherImpactTanamanRusak  WeatherImpact = "TANAMAN_RUSAK"
    WeatherImpactLainnya       WeatherImpact = "LAINNYA"
)


type MainConstraint string

const (
    MainConstraintIrigasiSulit MainConstraint = "IRIGASI_SULIT"
    MainConstraintHargaRendah  MainConstraint = "HARGA_RENDAH"
    MainConstraintPupuk        MainConstraint = "PUPUK"
    MainConstraintHama         MainConstraint = "HAMA"
    MainConstraintIklim        MainConstraint = "IKLIM"
    MainConstraintAksesPasar   MainConstraint = "AKSES_PASAR"
    MainConstraintLainnya      MainConstraint = "LAINNYA"
)


type FarmerHope string

const (
    FarmerHopeBantuanAlsintan FarmerHope = "BANTUAN_ALSINTAN"
    FarmerHopeBibitPupuk      FarmerHope = "BIBIT_PUPUK"
    FarmerHopeHargaStabil     FarmerHope = "HARGA_STABIL"
    FarmerHopePelatihan       FarmerHope = "PELATIHAN"
    FarmerHopeColdStorage     FarmerHope = "COLD_STORAGE"
    FarmerHopeLainnya         FarmerHope = "LAINNYA"
)


type TrainingNeeded string

const (
    TrainingNeededPHT              TrainingNeeded = "PHT"
    TrainingNeededPupukOrganik     TrainingNeeded = "PUPUK_ORGANIK"
    TrainingNeededPascapanen       TrainingNeeded = "PASCAPANEN"
    TrainingNeededPemasaranDigital TrainingNeeded = "PEMASARAN_DIGITAL"
    TrainingNeededGreenhouseIoT    TrainingNeeded = "GREENHOUSE_IOT"
    TrainingNeededLainnya          TrainingNeeded = "LAINNYA"
)


type UrgentNeeds string

const (
    UrgentNeedsPerbaikanIrigasi UrgentNeeds = "PERBAIKAN_IRIGASI"
    UrgentNeedsBibitPupukSegera UrgentNeeds = "BIBIT_PUPUK_SEGERA"
    UrgentNeedsReplanting       UrgentNeeds = "REPLANTING"
    UrgentNeedsColdStorage      UrgentNeeds = "COLD_STORAGE"
    UrgentNeedsLainnya          UrgentNeeds = "LAINNYA"
)


type WaterAccess string

const (
    WaterAccessMudah   WaterAccess = "MUDAH"
    WaterAccessJauh    WaterAccess = "JAUH"
    WaterAccessTidakAda WaterAccess = "TIDAK_ADA"
    WaterAccessLainnya WaterAccess = "LAINNYA"
)