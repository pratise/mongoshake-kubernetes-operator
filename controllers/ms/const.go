package ms

const (
	MongoshakeConfigVolClaimName = "mongoshake-collector"
	MongoshakeConfigDir          = "/mongoshake/config"

	timezone    = "timezone"
	timezoneDir = "/etc/localtime"
	WorkingDir  = "/mongoshake"

	MongoshakeLogVolClaimName = "ms-log"
	MongoshakeConfigName      = "ms-cfg"
	MongoshakeLogDir          = "/mongoshake/log"
	// 监控端口
	FullSyncPortName = "fullsync" // 全量同步监控端口
	IncrSyncPortName = "incrsync"
)
