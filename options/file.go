package options

import "gopkg.in/natefinch/lumberjack.v2"

// fileOption 日志输出到文件配置
// path 存储路径
// maxSize 在切割前。日志文件最大大小
// maxBackups 日志最多保留的文件个数
// maxAges 日志最多保留的天数
// compress 是否压缩
// localTime 时间是否本地化
type fileOption struct {
	path       string
	maxSize    int
	maxBackups int
	maxAges    int
	compress   bool
	localTime  bool
}

func lumberjackLogger(config Options) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.file.path,
		MaxSize:    config.file.maxSize,
		MaxAge:     config.file.maxAges,
		MaxBackups: config.file.maxBackups,
		LocalTime:  config.file.localTime,
		Compress:   config.file.compress,
	}
}
