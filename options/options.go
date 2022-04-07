package options

import (
	"context"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Type int

const (
	Console Type = iota
	Elastic
	File
)

type Option func(*Options)

var defaultConfig = Options{
	AppName: "log",
	Level:   zap.DebugLevel,
	LogType: Console,
	file: fileOption{
		path:       "./logs",
		maxSize:    30, //nolint:gomnd
		maxBackups: 10, //nolint:gomnd
		maxAges:    10, //nolint:gomnd
		compress:   true,
		localTime:  true,
	},
	ctx: context.Background(),
}

// Options 日志配置
type Options struct {
	Level   zapcore.Level
	AppName string
	LogType Type
	file    fileOption
	out     io.Writer
	ctx     context.Context
}

func NewOptions() *Options {
	return &Options{}
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

func Out(o io.Writer) Option {
	return func(options *Options) {
		options.out = o
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
