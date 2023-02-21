package navigations

import (
	"fyne.io/fyne/v2"
	"time"
	"window_handler/config"
)

const isDev = true

var timeCycleMap = make(map[string]time.Duration)
var dayCycleMap = make(map[string]time.Weekday)
var dayArrayList = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

var classicSize = fyne.Size{
	Height: config.WindowHeight * 0.8,
	Width:  config.WindowWidth,
}
var disableRootCache = make(map[string]map[fyne.Disableable]disableRoot)

/*
*
Selectivity between some components presents strong correlation
*/
type disableRoot struct {
	child []fyne.Disableable
}

func (p *disableRoot) addChild(childs ...fyne.Disableable) {
	for _, v := range childs {
		p.child = append(p.child, v)
	}
}

func (p *disableRoot) disableChild() {
	for _, v := range p.child {
		v.Disable()
	}
}

func (p *disableRoot) enableChild() {
	for _, v := range p.child {
		v.Enable()
	}
}

type Navigation struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
}

var (
	Navigations = map[string]Navigation{
		"localSync": {
			"Local Sync",
			"",
			GetBatchLocalSyncComponent,
			true,
		},
		"localBatchSync": {
			"Local Batch Sync",
			"Click start button to begin sync",
			GetBatchLocalSyncComponent,
			true,
		},
		"localSingleSync": {
			"Local Single Sync",
			"Click start button to begin sync.",
			GetSingleLocalSyncComponent,
			true,
		},
		"remoteSync": {
			"Remote",
			"QNQ Target Info",
			GetRemoteSyncComponent,
			true,
		},
		"remoteSingleSync": {
			"Remote sync",
			"Remote",
			GetRemoteSingleComponent,
			true,
		},
		"systemInfo": {
			"System Information",
			"",
			GetLocalSystemInfoComponent,
			true,
		},
		"diskInfo": {
			"Disk Information",
			"Basic Disk Information",
			GetDiskInfoComponent,
			true,
		},
	}
	//设置菜单树
	NavigationIndex = map[string][]string{
		"":           {"localSync", "systemInfo", "remoteSync"},
		"localSync":  {"localBatchSync", "localSingleSync"},
		"systemInfo": {"diskInfo"},
		"remoteSync": {"remoteSingleSync"},
	}
)

func initTimeCycle() {
	timeCycleMap["Second"] = time.Second
	timeCycleMap["Minute"] = time.Minute
	timeCycleMap["Hour"] = time.Hour

	dayCycleMap[dayArrayList[0]] = time.Sunday
	dayCycleMap[dayArrayList[1]] = time.Monday
	dayCycleMap[dayArrayList[2]] = time.Tuesday
	dayCycleMap[dayArrayList[3]] = time.Wednesday
	dayCycleMap[dayArrayList[4]] = time.Thursday
	dayCycleMap[dayArrayList[5]] = time.Friday
	dayCycleMap[dayArrayList[6]] = time.Saturday
}
