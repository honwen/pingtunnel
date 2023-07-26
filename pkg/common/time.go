package common

import (
	"time"
)

var gnowsecond time.Time

func init() {
	gnowsecond = time.Now()
	go updateNowInSecond()
}

func GetNowUpdateInSecond() time.Time {
	return gnowsecond
}

func updateNowInSecond() {
	defer CrashLog()

	for {
		gnowsecond = time.Now()
		Sleep(1)
	}
}

func Sleep(sec int) {
	last := time.Now()
	for time.Since(last) < time.Second*time.Duration(sec) {
		time.Sleep(time.Millisecond * 100)
	}
}
