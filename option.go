package scache

type OptionF func(*option)

type option struct {
	AutoClean         bool
	CleanIntervalTime int // 默认扫描间隔
}

// 开启监管者,定时间隔gap
func WithAutoClean(gap int) OptionF {
	return func(opt *option) {
		opt.AutoClean = true
		opt.CleanIntervalTime = gap
	}
}
