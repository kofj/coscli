package util

const (
	Standard              = "STANDARD"
	StandardIA            = "STANDARD_IA"
	IntelligentTiering    = "INTELLIGENT_TIERING"
	Archive               = "ARCHIVE"
	DeepArchive           = "DEEP_ARCHIVE"
	MAZStandard           = "MAZ_STANDARD"
	MAZStandardIA         = "MAZ_STANDARD_IA"
	MAZIntelligentTiering = "MAZ_INTELLIGENT_TIERING"
	MAZArchive            = "MAZ_ARCHIVE"
)

const (
	CommandCP      = "cp"
	CommandSync    = "sync"
	CommandLs      = "ls"
	CommandRm      = "rm"
	CommandRestore = "restore"
)

const (
	TypeSnapshotPath   = "snapshotPath"
	TypeFailOutputPath = "failOutputPath"
)

const (
	Version             string = "v1.0.6"
	Package             string = "coscli"
	SchemePrefix        string = "cos://"
	CosSeparator        string = "/"
	IncludePrompt              = "--include"
	ExcludePrompt              = "--exclude"
	ChannelSize         int    = 1000
	MaxSyncNumbers             = 1000000
	MaxDeleteBatchCount int    = 1000
	SnapshotConnector          = "==>"
	OfsMaxRenderNum     int    = 100
)

const (
	CpTypeUpload CpType = iota
	CpTypeDownload
	CpTypeCopy
)

const (
	DU_TYPE_TOTAL          = 1
	DU_TYPE_CATEGORIZATION = 2
)

// 版本控制状态
const (
	VersionStatusSuspended = "Suspended"
	VersionStatusEnabled   = "Enabled"
)
