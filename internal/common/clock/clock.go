package clock

import "time"

// Clock 定義可抽換的時間控制介面
// 適用於需要可測試、模擬時間的場景
// Sleep 和 After 實作預設為非精確模擬（不使用真實時間）
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	After(d time.Duration) <-chan time.Time
}

// RealClock 使用 time 標準函式實作 Clock，適用於 production 環境
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

func (RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
