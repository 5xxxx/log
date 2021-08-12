package options

import (
	"context"
	"io"
	"os"
	"time"

	redis2 "github.com/5xxxx/log/redis"

	"go.elastic.co/ecszap"

	"github.com/go-redis/redis/v8"
	"github.com/mattn/go-colorable" //nolint:gci
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Type int

const (
	Term Type = iota
	Elastic
	File
)

type Option func(*Options)

var defaultConfig = Options{
	AppName: "log",
	Level:   zap.DebugLevel,
	LogType: Term,
	file: fileOption{
		path:       "./logs",
		maxSize:    30, //nolint:gomnd
		maxBackups: 10, //nolint:gomnd
		maxAges:    10, //nolint:gomnd
		compress:   true,
		localTime:  true,
	},
	redis: redisOption{
		redis:       nil,
		redisLogKey: "app.log",
	},
	ctx: context.Background(),
}

// Options 日志配置
type Options struct {
	Level   zapcore.Level
	AppName string
	LogType Type
	file    fileOption
	redis   redisOption
	ctx     context.Context
}

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

// redisOption 日志输出到redis配置
type redisOption struct {
	redis       *redis.Client
	redisLogKey string
}

func NewOptions() *Options {
	return &Options{}
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

func Leve(l zapcore.Level) Option {
	return func(options *Options) {
		options.Level = l
	}
}

func LogType(t Type) Option {
	return func(options *Options) {
		options.LogType = t
	}
}

func AppName(app string) Option {
	return func(options *Options) {
		options.AppName = app
	}
}

func Path(p string) Option {
	return func(options *Options) {
		options.file.path = p
	}
}

func MaxBackups(size int) Option {
	return func(options *Options) {
		options.file.maxBackups = size
	}
}

func MaxSize(size int) Option {
	return func(options *Options) {
		options.file.maxSize = size
	}
}

func MaxAges(size int) Option {
	return func(options *Options) {
		options.file.maxAges = size
	}
}

func Compress(c bool) Option {
	return func(options *Options) {
		options.file.compress = c
	}
}

func LocalTime(l bool) Option {
	return func(options *Options) {
		options.file.localTime = l
	}
}

func RedisClient(r *redis.Client) Option {
	return func(options *Options) {
		options.redis.redis = r
	}
}

func RedisLogKey(r string) Option {
	return func(options *Options) {
		options.redis.redisLogKey = r
	}
}

func GetOptions(opts ...Option) Options {
	conf := defaultConfig
	for _, op := range opts {
		op(&conf)
	}
	return conf
}

func (option Options) Build() *zap.Logger {
	skip := 2
	if option.LogType == Elastic {
		skip = 2
	}
	return zap.New(ZapCore(option), zap.AddCaller(), zap.AddCallerSkip(skip))
}

func ZapCore(option Options) zapcore.Core {
	switch option.LogType {
	case Term:
		return TermCore(option)
	case Elastic:
		return ElasticCore(option)
	case File:
		return FileCore(option)
	}
	return nil
}

func FileCore(option Options) zapcore.Core {
	cfgConsole := zapcore.EncoderConfig{
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("2006-01-02 15:04:05")
		},
		TimeKey:        "time",
		LevelKey:       "Level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(option.Level)

	return zapcore.NewCore(zapcore.NewConsoleEncoder(cfgConsole),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			zapcore.AddSync(lumberjackLogger(option))),
		atomicLevel)
}

func TermCore(option Options) zapcore.Core {
	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder
	aa.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("2006-01-02 15:04:05")
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		option.Level,
	)
}

func ElasticCore(option Options) zapcore.Core {
	if option.writer() == nil {
		panic("writer can't be nil")
	}

	var topicErrors = zapcore.AddSync(option.writer()) //nolint:gofumpt
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	return ecszap.NewCore(encoderConfig, topicErrors, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return option.Level <= level
	}))
}

func (option Options) writer() io.Writer {
	if option.redis.redis != nil {
		return redis2.NewRedisWriter(option.ctx, option.redis.redisLogKey, option.redis.redis)
	}
	return nil
}
