package log

import (
	"fmt"

	"github.com/5xxxx/log/options"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg interface{}, fields ...zap.Field)
	Info(msg interface{}, fields ...zap.Field)
	Warn(msg interface{}, fields ...zap.Field)
	Error(msg interface{}, fields ...zap.Field)
	Panic(msg interface{}, fields ...zap.Field)
}

type logs struct {
	logger *zap.Logger
	config options.Options
}

var logger Logger

func Init(opts ...options.Option) {
	logger = NewLog(opts...)
}

func NewLog(opts ...options.Option) Logger {
	c := options.GetOptions(opts...)
	return &logs{
		logger: c.Build(),
		config: c,
	}
}

func (l *logs) configFields(fields ...zap.Field) []zap.Field {
	if l.config.LogType == options.Elastic {
		fields = append(fields, zap.String("app_name", l.config.AppName))
	}
	return fields
}

func (l *logs) Debug(msg interface{}, fields ...zap.Field) {
	l.logger.Debug(fmt.Sprint(msg), l.configFields(fields...)...)
}

func (l *logs) Info(msg interface{}, fields ...zap.Field) {
	l.logger.Info(fmt.Sprint(msg), l.configFields(fields...)...)
}

func (l *logs) Warn(msg interface{}, fields ...zap.Field) {
	l.logger.Warn(fmt.Sprint(msg), l.configFields(fields...)...)
}

func (l *logs) Error(msg interface{}, fields ...zap.Field) {
	l.logger.Error(fmt.Sprint(msg), l.configFields(fields...)...)
}

func (l *logs) Panic(msg interface{}, fields ...zap.Field) {
	l.logger.Panic(fmt.Sprint(msg), l.configFields(fields...)...)
}

func Debug(msg interface{}, fields ...zap.Field) {
	logger.Debug(fmt.Sprint(msg), fields...)
}

func Info(msg interface{}, fields ...zap.Field) {
	logger.Info(fmt.Sprint(msg), fields...)
}

func Warn(msg interface{}, fields ...zap.Field) {
	logger.Warn(fmt.Sprint(msg), fields...)
}

func Error(msg interface{}, fields ...zap.Field) {
	logger.Error(fmt.Sprint(msg), fields...)
}

func Panic(msg interface{}, fields ...zap.Field) {
	logger.Panic(fmt.Sprint(msg), fields...)
}
