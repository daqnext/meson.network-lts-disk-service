package downloadtaskmgr

type DownloadInfo struct {
	TargetUrl    string
	BindName     string
	FileName     string
	Continent    string
	Country      string
	Area         string
	SavePath     string
	DownloadType string
	OriginRegion string
	TargetRegion string
}

type TaskStatus string

const TaskUnStart TaskStatus = "unstart"
const TaskBreak TaskStatus = "break"
const TaskDownloading TaskStatus = "downloading"

type DownloadTask struct {
	DownloadInfo
	Id              uint64
	Status          TaskStatus
	FileSize        int64
	SpeedKBs        float64
	DownloadedSize  int64
	TryTimes        int
	StartTime       int64
	ZeroSpeedSec    int
	DownloadChannel *DownloadChannel
}
