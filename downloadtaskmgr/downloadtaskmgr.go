package downloadtaskmgr

import "math"

type DownloadTaskMgr struct {
	DownloadChannelArray []*DownloadChannel
}

type DownloadChannel struct {
	SpeedLimitKBs           int64
	CountLimit              int
	RunningCountControlChan chan bool
	IdleChan                chan *DownloadTask
}

const NewRunningTaskCount = 7

var newRunningTaskControlChan = make(chan bool, NewRunningTaskCount)

func New() *DownloadTaskMgr {
	taskMgr := &DownloadTaskMgr{}
	taskMgr.DownloadChannelArray = []*DownloadChannel{
		{SpeedLimitKBs: 200, CountLimit: 4, RunningCountControlChan: make(chan bool, 4), IdleChan: make(chan *DownloadTask, 1024*5)},           //0-200KB/s
		{SpeedLimitKBs: 1500, CountLimit: 3, RunningCountControlChan: make(chan bool, 3), IdleChan: make(chan *DownloadTask, 1024*5)},          //100-1500KB/s
		{SpeedLimitKBs: math.MaxInt64, CountLimit: 3, RunningCountControlChan: make(chan bool, 3), IdleChan: make(chan *DownloadTask, 1024*3)}, //>1500KB/s
	}
	return taskMgr
}
