
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
    StatusKerusakan ReportStatusType = "KERUSAKAN"
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

type FarmerGroupType string

const (
    FarmerGroupTypePoktan       FarmerGroupType = "POKTAN"
    FarmerGroupTypeGapoktan     FarmerGroupType = "GAPOKTAN"
)

const (
    WorkTypePerbaikanAtap       WorkType = "PERBAIKAN_ATAP"
    WorkTypePerbaikanDinding    WorkType = "PERBAIKAN_DINDING"
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
    InstitutionUnitUptIrigasi     InstitutionUnitType = "UPT_IRIGASI"
    InstitutionUnitPoktan      InstitutionUnitType = "POKTAN"
    InstitutionUnitDinasPupr InstitutionUnitType = "DINAS_PUPR"
    InstitutionUnitUptJalan InstitutionUnitType = "UPT_JALAN"
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

// type RoadType string

// const (
//     RoadTypeJalanNasional   RoadType = "JALAN_NASIONAL"
//     RoadTypeJalanProvinsi   RoadType = "JALAN_PROVINSI"
//     RoadTypeJalanKabupaten  RoadType = "JALAN_KABUPATEN"
//     RoadTypeJalanKota       RoadType = "JALAN_KOTA"
//     RoadTypeJalanDesa       RoadType = "JALAN_DESA"
//     RoadTypeJalanLingkungan RoadType = "JALAN_LINGKUNGAN"
// )


// type RoadClass string

// const (
//     RoadClassArteri      RoadClass = "ARTERI"      
//     RoadClassKolektor    RoadClass = "KOLEKTOR"    
//     RoadClassLokal       RoadClass = "LOKAL"       
//     RoadClassLingkungan  RoadClass = "LINGKUNGAN"  
// )

type RoadDamageLevel string

const (
    RoadDamageLevelMinor    RoadDamageLevel = "RINGAN"   
    RoadDamageLevelModerate RoadDamageLevel = "SEDANG"   
    RoadDamageLevelSevere   RoadDamageLevel = "BERAT"    
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


const (
    LandStatusOwned  LandStatus = "MILIK_SENDIRI"
    LandStatusRented LandStatus = "SEWA"
    LandStatusFree   LandStatus = "PINJAM_BEBAS_SEWA"
    LandStatusOther  LandStatus = "LAINNYA"
)

type PavementType string

const (
    PavementAspalFlexible PavementType = "ASPAL_FLEXIBLE"
    PavementBetonRigid    PavementType = "BETON_RIGID"
    PavementPaving        PavementType = "PAVING"
    PavementJalanTanah    PavementType = "JALAN_TANAH"
)


type RoadDamageType string

const (
    RoadDamageLubang            RoadDamageType = "LUBANG_POTHOLES"
    RoadDamageRetakBuaya        RoadDamageType = "RETAK_KULIT_BUAYA"
    RoadDamageAmblas            RoadDamageType = "AMBLAS_LONGSOR"
    RoadDamagePermukaanAus      RoadDamageType = "PERMUKAAN_AUS_RAVELING"
    RoadDamageGenaganDrainase   RoadDamageType = "GENANGAN_AIR_DRAINASE_BURUK"
    RoadDamageRetakMemanjang    RoadDamageType = "RETAK_MEMANJANG"
    RoadDamageRetakMelintang    RoadDamageType = "RETAK_MELINTANG"
    RoadDamageRetakBlok         RoadDamageType = "RETAK_BLOK"
    RoadDamageGelombang         RoadDamageType = "GELOMBANG"
    RoadDamageTepiJalan         RoadDamageType = "KERUSAKAN_TEPI"
    RoadDamageDrainase          RoadDamageType = "KERUSAKAN_DRAINASE"
    RoadDamageJembatan          RoadDamageType = "KERUSAKAN_JEMBATAN"
    RoadDamagePerlengkapan      RoadDamageType = "KERUSAKAN_PERLENGKAPAN"
    RoadDamageLainnya           RoadDamageType = "LAINNYA"
)


type BridgeStructureType string

const (
    BridgeStructureBetonBertulang BridgeStructureType = "BETON_BERTULANG"
    BridgeStructureBaja           BridgeStructureType = "BAJA"
    BridgeStructureKayu           BridgeStructureType = "KAYU"
)


type BridgeDamageType string

const (
    BridgeDamageLantaiRetak    BridgeDamageType = "LANTAI_JEMBATAN_RETAK_RUSAK"
    BridgeDamageOpritAmblas    BridgeDamageType = "OPRIT_ABUTMENT_AMBLAS"
    BridgeDamageRangkaRetak    BridgeDamageType = "RANGKA_UTAMA_RETAK"
    BridgeDamagePondasiTerseret BridgeDamageType = "PONDASI_TERSERET_ARUS"
    BridgeDamageLainnya        BridgeDamageType = "LAINNYA"
)


type BridgeDamageLevel string

const (
    BridgeDamageLevelRingan      BridgeDamageLevel = "RINGAN"
    BridgeDamageLevelSedang      BridgeDamageLevel = "SEDANG"
    BridgeDamageLevelSevere      BridgeDamageLevel = "BERAT_TIDAK_LAYAK"
)


type TrafficCondition string

const (
    TrafficConditionNormal         TrafficCondition = "MASIH_BISA_DILALUI"
    TrafficConditionOneLane        TrafficCondition = "HANYA_SATU_LAJUR_BISA_DILALUI"
    TrafficConditionBlocked        TrafficCondition = "TIDAK_BISA_DILALUI_PUTUS"
)


type RoadUrgencyLevel string

const (
    RoadUrgencyEmergency RoadUrgencyLevel = "DARURAT"   
    RoadUrgencyHigh      RoadUrgencyLevel = "CEPAT"     
    RoadUrgencyMedium    RoadUrgencyLevel = "RUTIN"     
    RoadUrgencyLow       RoadUrgencyLevel = "RENDAH"    
)


type TrafficImpact string

const (
    TrafficImpactMinimal         TrafficImpact = "MINIMAL"           
    TrafficImpactReduced         TrafficImpact = "TERGANGGU"         
    TrafficImpactSeverelyReduced TrafficImpact = "SANGAT_TERGANGGU"  
    TrafficImpactBlocked         TrafficImpact = "TERPUTUS"          
)

// Enhanced Food Commodities
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

// Enhanced Horticulture Commodities
type HorticultureCommodity string

const (
    HortiCommoditySayuran      HorticultureCommodity = "SAYURAN"
    HortiCommodityBuah         HorticultureCommodity = "BUAH"
    HortiCommodityFlorikultura HorticultureCommodity = "FLORIKULTURA"
    HortiCommodityObatTradsional HorticultureCommodity = "TANAMAN_OBAT_TRADISIONAL"
)

// Enhanced Plantation Commodities
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
    PlantationCommodityVanili     PlantationCommodity = "VANILI"
    PlantationCommodityLada       PlantationCommodity = "LADA"
    PlantationCommodityPala       PlantationCommodity = "PALA"
    PlantationCommodityLainnya    PlantationCommodity = "LAINNYA"
)

// Enhanced Land Status
type LandStatus string

const (
    LandStatusMilikSendiri LandStatus = "MILIK_SENDIRI"
    LandStatusSewa         LandStatus = "SEWA"
    LandStatusBagiHasil    LandStatus = "BAGI_HASIL"
    LandStatusBebasSewaOff LandStatus = "PINJAM_BEBAS_SEWA"
    LandStatusHibah        LandStatus = "HIBAH"
    LandStatusLainnya      LandStatus = "LAINNYA"
)

// Enhanced Growth Phases for Food Crops
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
    GrowthPhasePascaPanen     GrowthPhase = "PASCA_PANEN"
    GrowthPhaseLainnya        GrowthPhase = "LAINNYA"
)

// Enhanced Horticulture Growth Phases
type HortiGrowthPhase string

const (
    HortiGrowthPhasePersemaian   HortiGrowthPhase = "PERSEMAIAN"
    HortiGrowthPhasePembibitan   HortiGrowthPhase = "PEMBIBITAN"
    HortiGrowthPhaseTanam        HortiGrowthPhase = "TANAM"
    HortiGrowthPhaseVegetatif    HortiGrowthPhase = "VEGETATIF"
    HortiGrowthPhasePembungaan   HortiGrowthPhase = "PEMBUNGAAN"
    HortiGrowthPhasePembuahan    HortiGrowthPhase = "PEMBUAHAN"
    HortiGrowthPhasePanen        HortiGrowthPhase = "PANEN"
    HortiGrowthPhasePascaPanen   HortiGrowthPhase = "PASCA_PANEN"
    HortiGrowthPhaseLainnya      HortiGrowthPhase = "LAINNYA"
)

// Enhanced Plantation Growth Phases
type PlantationGrowthPhase string

const (
    PlantationGrowthPhaseBibit             PlantationGrowthPhase = "BIBIT_PERSEMAIAN"
    PlantationGrowthPhaseTanamanMuda       PlantationGrowthPhase = "TANAMAN_MUDA_TBM"
    PlantationGrowthPhaseTanamanMenghasilkan PlantationGrowthPhase = "TANAMAN_MENGHASILKAN_TM"
    PlantationGrowthPhaseReplanting        PlantationGrowthPhase = "REPLANTING"
    PlantationGrowthPhasePanen             PlantationGrowthPhase = "PANEN"
    PlantationGrowthPhasePemerliharaan     PlantationGrowthPhase = "PEMELIHARAAN"
    PlantationGrowthPhaseLainnya           PlantationGrowthPhase = "LAINNYA"
)

// Enhanced Delay Reasons
type DelayReason string

const (
    DelayReasonTidakAda       DelayReason = "TIDAK_ADA"
    DelayReasonHujanTerus     DelayReason = "HUJAN_TERUS_MENERUS"
    DelayReasonKekeringan     DelayReason = "KEKERINGAN"
    DelayReasonBibitTerlambat DelayReason = "BIBIT_TERLAMBAT"
    DelayReasonBanjir         DelayReason = "BANJIR"
    DelayReasonSerangan       DelayReason = "SERANGAN_HAMA_PENYAKIT"
    DelayReasonPermasalahanModal DelayReason = "PERMASALAHAN_MODAL"
    DelayReasonTenagaKerja    DelayReason = "KETERBATASAN_TENAGA_KERJA"
    DelayReasonLainnya        DelayReason = "LAINNYA"
)

// Enhanced Technology Methods for Food Crops
type TechnologyMethod string

const (
    TechnologyMethodTidakAda        TechnologyMethod = "TIDAK_ADA"
    TechnologyMethodJajarLegowo     TechnologyMethod = "JAJAR_LEGOWO"
    TechnologyMethodDroneSemprot    TechnologyMethod = "DRONE_SEMPROT"
    TechnologyMethodPupukOrganik    TechnologyMethod = "PUPUK_ORGANIK_HAYATI"
    TechnologyMethodIrigasiPompa    TechnologyMethod = "IRIGASI_POMPA"
    TechnologyMethodBibitUnggul     TechnologyMethod = "VARIETAS_UNGGUL_BERSERTIFIKAT"
    TechnologyMethodPengolahan      TechnologyMethod = "PENGOLAHAN_TANAH_MINIMUM"
    TechnologyMethodPengendalianHama TechnologyMethod = "PENGENDALIAN_HAMA_TERPADU"
    TechnologyMethodLainnya         TechnologyMethod = "LAINNYA"
)

// Enhanced Horticulture Technology
type HortiTechnology string

const (
    HortiTechnologyTidakAda         HortiTechnology = "TIDAK_ADA"
    HortiTechnologyGreenhouse       HortiTechnology = "GREENHOUSE_SCREEN_HOUSE"
    HortiTechnologyMulsaPlastik     HortiTechnology = "MULSA_PLASTIK"
    HortiTechnologyIrigasiTetes     HortiTechnology = "IRIGASI_TETES_SPRINKLER"
    HortiTechnologyDroneSemprot     HortiTechnology = "DRONE_SEMPROT"
    HortiTechnologySensorIoT        HortiTechnology = "SENSOR_IOT_KELEMBAPAN"
    HortiTechnologyPupukOrganik     HortiTechnology = "PUPUK_ORGANIK_HAYATI"
    HortiTechnologyHydroponik       HortiTechnology = "HIDROPONIK_AEROPONIK"
    HortiTechnologyVerticalFarming  HortiTechnology = "VERTICAL_FARMING"
    HortiTechnologyBibitUnggul      HortiTechnology = "BIBIT_VARIETAS_UNGGUL"
    HortiTechnologyLainnya          HortiTechnology = "LAINNYA"
)

// Enhanced Plantation Technology
type PlantationTechnology string

const (
    PlantationTechnologyTidakAda        PlantationTechnology = "TIDAK_ADA"
    PlantationTechnologyPeremajaan      PlantationTechnology = "PEREMAJAAN_TANAMAN_REPLANTING"
    PlantationTechnologyPupukOrganik    PlantationTechnology = "PUPUK_ORGANIK_HAYATI"
    PlantationTechnologyIrigasiTetes    PlantationTechnology = "IRIGASI_TETES_SPRINKLER"
    PlantationTechnologyDroneMonitoring PlantationTechnology = "DRONE_MONITORING"
    PlantationTechnologyAgroforestry    PlantationTechnology = "SISTEM_AGROFORESTRY"
    PlantationTechnologyPengeringanModern PlantationTechnology = "TEKNOLOGI_PENGERINGAN_MODERN"
    PlantationTechnologyFermentasi      PlantationTechnology = "TEKNIK_FERMENTASI_TERKONTROL"
    PlantationTechnologyLainnya         PlantationTechnology = "LAINNYA"
)

// Enhanced Post Harvest Problems
type PostHarvestProblem string

const (
    PostHarvestProblemTidakAda           PostHarvestProblem = "TIDAK_ADA"
    PostHarvestProblemSusutTinggi        PostHarvestProblem = "SUSUT_TINGGI_HASIL_CEPAT_BUSUK"
    PostHarvestProblemKeterbatasanGudang PostHarvestProblem = "KETERBATASAN_GUDANG_PENDINGIN_COLD_STORAGE"
    PostHarvestProblemKesulitanKemasan   PostHarvestProblem = "KESULITAN_KEMASAN_GRADING"
    PostHarvestProblemKesulitanAkses     PostHarvestProblem = "KESULITAN_AKSES_PASAR"
    PostHarvestProblemHargaRendah        PostHarvestProblem = "HARGA_JUAL_RENDAH"
    PostHarvestProblemTengkulak          PostHarvestProblem = "MONOPOLI_TENGKULAK"
    PostHarvestProblemTransportasi       PostHarvestProblem = "KETERBATASAN_TRANSPORTASI"
    PostHarvestProblemLainnya            PostHarvestProblem = "LAINNYA"
)

// Enhanced Production Problems
type ProductionProblem string

const (
    ProductionProblemRendahnyaProduktivitas ProductionProblem = "RENDAHNYA_PRODUKTIVITAS"
    ProductionProblemHargaJualFluktuatif    ProductionProblem = "HARGA_JUAL_FLUKTUATIF"
    ProductionProblemSeranganHama           ProductionProblem = "SERANGAN_HAMA_PENYAKIT"
    ProductionProblemPerluReplanting        ProductionProblem = "PERLU_REPLANTING"
    ProductionProblemKekuranganModal        ProductionProblem = "KEKURANGAN_MODAL_USAHA"
    ProductionProblemKeterbatasanLahan      ProductionProblem = "KETERBATASAN_LAHAN"
    ProductionProblemKualitasRendah         ProductionProblem = "KUALITAS_HASIL_RENDAH"
    ProductionProblemTenagaKerja            ProductionProblem = "KETERBATASAN_TENAGA_KERJA"
    ProductionProblemLainnya                ProductionProblem = "LAINNYA"
)

// Enhanced Pest Disease Types with Commodity Specific
type PestDiseaseType string

const (
    // Food Crops Pest & Disease
    PestDiseaseUlatGrayak     PestDiseaseType = "ULAT_GRAYAK"
    PestDiseaseWerengCoklat   PestDiseaseType = "WERENG_COKLAT"
    PestDiseaseTikus          PestDiseaseType = "TIKUS"
    PestDiseaseBusukDaun      PestDiseaseType = "BUSUK_DAUN"
    PestDiseaseBLAST          PestDiseaseType = "BLAST_PADI"
    PestDiseaseBusukBatang    PestDiseaseType = "BUSUK_BATANG"
    
    // Horticulture Pest & Disease
    PestDiseaseTrips          PestDiseaseType = "TRIPS"
    PestDiseaseLalatBuah      PestDiseaseType = "LALAT_BUAH"
    PestDiseaseAntraknosa     PestDiseaseType = "ANTRAKNOSA"
    PestDiseaseLayuFusarium   PestDiseaseType = "LAYU_FUSARIUM"
    PestDiseaseKutuDaun       PestDiseaseType = "KUTU_DAUN"
    PestDiseaseMosaik         PestDiseaseType = "VIRUS_MOSAIK"
    
    // Plantation Pest & Disease
    PestDiseasePBKo           PestDiseaseType = "PBKO_PENGGEREK_BUAH_KAKAO"
    PestDiseaseKaratDaun      PestDiseaseType = "KARAT_DAUN"
    PestDiseaseHamaBorrer     PestDiseaseType = "HAMA_BORER"
    PestDiseaseHamaTikus      PestDiseaseType = "HAMA_TIKUS"
    PestDiseasePenyakitAkar   PestDiseaseType = "PENYAKIT_AKAR"
    
    PestDiseaseLainnya        PestDiseaseType = "LAINNYA"
)

// Enhanced Affected Area Levels
type AffectedAreaLevel string

const (
    AffectedAreaLevelKurang10   AffectedAreaLevel = "KURANG_10_PERSEN"
    AffectedAreaLevel10Sampai25 AffectedAreaLevel = "10_SAMPAI_25_PERSEN"
    AffectedAreaLevel25Sampai50 AffectedAreaLevel = "25_SAMPAI_50_PERSEN"
    AffectedAreaLevelLebih50    AffectedAreaLevel = "LEBIH_50_PERSEN"
    AffectedAreaLevelSeluruh    AffectedAreaLevel = "SELURUH_AREA"
    AffectedAreaLevelLainnya    AffectedAreaLevel = "LAINNYA"
)

// Enhanced Control Actions
type ControlAction string

const (
    ControlActionBelumDitangani     ControlAction = "BELUM_DITANGANI"
    ControlActionSemprotInsektisida ControlAction = "SEMPROT_INSEKTISIDA_KIMIA"
    ControlActionSemprotBiopestisida ControlAction = "SEMPROT_BIOPESTISIDA"
    ControlActionPasangPerangkap    ControlAction = "PASANG_PERANGKAP_FEROMON"
    ControlActionSanitasiKebun      ControlAction = "SANITASI_KEBUN_ROTASI_TANAMAN"
    ControlActionAgenHayati         ControlAction = "PELEPASAN_AGEN_HAYATI"
    ControlActionVarietasTahan      ControlAction = "TANAM_VARIETAS_TAHAN"
    ControlActionPHT                ControlAction = "PENERAPAN_PHT"
    ControlActionLainnya            ControlAction = "LAINNYA"
)

// Enhanced Weather Conditions
type WeatherCondition string

const (
    WeatherConditionHujan        WeatherCondition = "HUJAN"
    WeatherConditionCerah        WeatherCondition = "CERAH"
    WeatherConditionMendung      WeatherCondition = "MENDUNG"
    WeatherConditionAnginKencang WeatherCondition = "ANGIN_KENCANG"
    WeatherConditionKekeringan   WeatherCondition = "KEKERINGAN"
    WeatherConditionBanjir       WeatherCondition = "BANJIR"
    WeatherConditionEkstrem      WeatherCondition = "CUACA_EKSTREM"
    WeatherConditionLainnya      WeatherCondition = "LAINNYA"
)

// Enhanced Weather Impact
type WeatherImpact string

const (
    WeatherImpactTidakAda        WeatherImpact = "TIDAK_ADA"
    WeatherImpactTanamanRebah    WeatherImpact = "TANAMAN_REBAH"
    WeatherImpactDaunMenguning   WeatherImpact = "DAUN_MENGUNING"
    WeatherImpactBuahRontok      WeatherImpact = "BUAH_RONTOK"
    WeatherImpactTanamanRusak    WeatherImpact = "TANAMAN_RUSAK"
    WeatherImpactGagalPanen      WeatherImpact = "GAGAL_PANEN"
    WeatherImpactTerlambatTanam  WeatherImpact = "TERLAMBAT_TANAM"
    WeatherImpactKekeringanLahan WeatherImpact = "KEKERINGAN_LAHAN"
    WeatherImpactLainnya         WeatherImpact = "LAINNYA"
)

// Enhanced Main Constraints
type MainConstraint string

const (
    MainConstraintIrigasiSulit MainConstraint = "IRIGASI_SULIT"
    MainConstraintHargaRendah  MainConstraint = "HARGA_RENDAH"
    MainConstraintPupuk        MainConstraint = "PUPUK_MAHAL_LANGKA"
    MainConstraintHama         MainConstraint = "SERANGAN_HAMA_PENYAKIT"
    MainConstraintIklim        MainConstraint = "PERUBAHAN_IKLIM"
    MainConstraintAksesPasar   MainConstraint = "AKSES_PASAR_TERBATAS"
    MainConstraintModal        MainConstraint = "KETERBATASAN_MODAL"
    MainConstraintTenagaKerja  MainConstraint = "TENAGA_KERJA_TERBATAS"
    MainConstraintTeknologi    MainConstraint = "TEKNOLOGI_TERBATAS"
    MainConstraintLahan        MainConstraint = "KETERBATASAN_LAHAN"
    MainConstraintLainnya      MainConstraint = "LAINNYA"
)

// Enhanced Farmer Hopes
type FarmerHope string

const (
    FarmerHopeBantuanAlsintan FarmerHope = "BANTUAN_ALSINTAN"
    FarmerHopeBibitPupuk      FarmerHope = "BANTUAN_BIBIT_PUPUK"
    FarmerHopeHargaStabil     FarmerHope = "HARGA_STABIL_TERJAMIN"
    FarmerHopePelatihan       FarmerHope = "PELATIHAN_TEKNOLOGI"
    FarmerHopeColdStorage     FarmerHope = "COLD_STORAGE_GUDANG"
    FarmerHopeAksesPasar      FarmerHope = "AKSES_PASAR_LANGSUNG"
    FarmerHopeBantuanModal    FarmerHope = "BANTUAN_MODAL_KREDIT"
    FarmerHopeIrigasi         FarmerHope = "PERBAIKAN_IRIGASI"
    FarmerHopeAsuransi        FarmerHope = "ASURANSI_PERTANIAN"
    FarmerHopeLainnya         FarmerHope = "LAINNYA"
)

// Enhanced Training Needs
type TrainingNeeded string

const (
    TrainingNeededPHT              TrainingNeeded = "PHT_PENGENDALIAN_HAMA_TERPADU"
    TrainingNeededPupukOrganik     TrainingNeeded = "PUPUK_ORGANIK_HAYATI"
    TrainingNeededPascapanen       TrainingNeeded = "TEKNOLOGI_PASCAPANEN"
    TrainingNeededPemasaranDigital TrainingNeeded = "PEMASARAN_DIGITAL"
    TrainingNeededGreenhouseIoT    TrainingNeeded = "GREENHOUSE_IOT"
    TrainingNeededBudidayaModern   TrainingNeeded = "BUDIDAYA_MODERN"
    TrainingNeededKeuanganUsaha    TrainingNeeded = "MANAJEMEN_KEUANGAN_USAHA"
    TrainingNeededKoperasi         TrainingNeeded = "MANAJEMEN_KOPERASI"
    TrainingNeededSertifikasi      TrainingNeeded = "SERTIFIKASI_ORGANIK"
    TrainingNeededLainnya          TrainingNeeded = "LAINNYA"
)

// Enhanced Urgent Needs
type UrgentNeeds string

const (
    UrgentNeedsPerbaikanIrigasi UrgentNeeds = "PERBAIKAN_IRIGASI"
    UrgentNeedsBibitPupukSegera UrgentNeeds = "BIBIT_PUPUK_SEGERA"
    UrgentNeedsReplanting       UrgentNeeds = "REPLANTING_SEGERA"
    UrgentNeedsColdStorage      UrgentNeeds = "COLD_STORAGE_MENDESAK"
    UrgentNeedsObatHama         UrgentNeeds = "OBAT_HAMA_PENYAKIT"
    UrgentNeedsModalDarurat     UrgentNeeds = "MODAL_DARURAT"
    UrgentNeedsAlsintan         UrgentNeeds = "ALAT_MESIN_PERTANIAN"
    UrgentNeedsPasarDarurat     UrgentNeeds = "AKSES_PASAR_DARURAT"
    UrgentNeedsLainnya          UrgentNeeds = "LAINNYA"
)

// Enhanced Water Access
type WaterAccess string

const (
    WaterAccessMudah     WaterAccess = "MUDAH_TERSEDIA"
    WaterAccessTerbatas  WaterAccess = "TERBATAS_MUSIMAN"
    WaterAccessJauh      WaterAccess = "JAUH_SULIT"
    WaterAccessTidakAda  WaterAccess = "TIDAK_ADA"
    WaterAccessBerbayar  WaterAccess = "TERSEDIA_BERBAYAR"
    WaterAccessLainnya   WaterAccess = "LAINNYA"
)

// Pest Disease Commodity Types
type PestDiseaseCommodityType string

const (
    PestDiseaseCommodityPangan      PestDiseaseCommodityType = "PANGAN"
    PestDiseaseCommodityHortikultura PestDiseaseCommodityType = "HORTIKULTURA"
    PestDiseaseCommodityPerkebunan  PestDiseaseCommodityType = "PERKEBUNAN"
)