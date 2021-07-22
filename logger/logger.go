package logger

import (
	"github.com/mkawserm/abesh/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
)

var (
	loggerOnce sync.Once
	loggerIns  *Factory
)

type SugaredLogger struct {
	*zap.SugaredLogger
}

func (l *SugaredLogger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

//func (l *SugaredLogger) SetLevel(_ logger.LogLevel) {
//}

type Factory struct {
	mMutex       sync.Mutex
	mPackageName string

	mZapConfig zap.Config
	mZapLogger *zap.Logger

	mLoggers        map[string]*zap.Logger
	mSugaredLoggers map[string]*zap.SugaredLogger
}

func (l *Factory) ChangeLogLevel(level string) {
	level = strings.TrimSpace(level)
	//fmt.Printf("'%s'",level)
	if level == "debug" {
		l.mZapConfig.Level.SetLevel(zap.DebugLevel)
	} else if level == "info" {
		l.mZapConfig.Level.SetLevel(zap.InfoLevel)
	} else if level == "warn" {
		l.mZapConfig.Level.SetLevel(zap.WarnLevel)
	} else if level == "error" {
		l.mZapConfig.Level.SetLevel(zap.ErrorLevel)
	} else if level == "panic" {
		l.mZapConfig.Level.SetLevel(zap.PanicLevel)
	} else if level == "fatal" {
		l.mZapConfig.Level.SetLevel(zap.FatalLevel)
	} else if level == "dpanic" {
		l.mZapConfig.Level.SetLevel(zap.DPanicLevel)
	} else {
		l.mZapConfig.Level.SetLevel(zap.ErrorLevel)
	}
}

func (l *Factory) SetupZapLogger(newLogger *zap.Logger) {
	l.mMutex.Lock()
	defer l.mMutex.Unlock()

	l.mZapLogger = newLogger
}

func (l *Factory) GetZapLogger() *zap.Logger {
	l.mMutex.Lock()
	defer l.mMutex.Unlock()

	if l.mZapLogger == nil {
		err := l.buildZapLogger()
		if err != nil {
			panic(err)
		}
	}

	if l.mLoggers == nil {
		l.mLoggers = make(map[string]*zap.Logger)
	}

	if l.mSugaredLoggers == nil {
		l.mSugaredLoggers = make(map[string]*zap.SugaredLogger)
	}

	return l.mZapLogger
}

func (l *Factory) S(pkgName string) *zap.SugaredLogger {
	l.mMutex.Lock()
	defer l.mMutex.Unlock()

	v, ok := l.mSugaredLoggers[pkgName]
	if ok {
		return v
	} else {
		sl := l.mZapLogger.Named(pkgName).Sugar()
		l.mSugaredLoggers[pkgName] = sl
		return sl
	}
}

func (l *Factory) L(pkgName string) *zap.Logger {
	l.mMutex.Lock()
	defer l.mMutex.Unlock()

	v, ok := l.mLoggers[pkgName]
	if ok {
		return v
	} else {
		nl := l.mZapLogger.Named(pkgName)
		l.mLoggers[pkgName] = nl
		return nl
	}
}

func (l *Factory) setup() {
	l.mLoggers = make(map[string]*zap.Logger)
	l.mSugaredLoggers = make(map[string]*zap.SugaredLogger)
	l.mZapConfig = newZapConfig()
}

func (l *Factory) buildZapLogger(options ...zap.Option) error {
	zl, err := l.mZapConfig.Build(options...)

	if err != nil {
		return err
	}

	l.mZapLogger = zl
	defer func() {
		_ = l.mZapLogger.Sync()
	}()

	return nil
}

// GetLoggerFactory get logger factory instance
func GetLoggerFactory() *Factory {
	return loggerIns
}

// CS returns sugared logger
func CS(pkgName string) *SugaredLogger {
	return &SugaredLogger{S(pkgName)}
}

// L returns zap logger
func L(pkgName string) *zap.Logger {
	return loggerIns.L(pkgName)
}

// S returns zap sugared logger
func S(pkgName string) *zap.SugaredLogger {
	return loggerIns.S(pkgName)
}

func newZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newZapConfig() zap.Config {
	level := strings.TrimSpace(conf.EnvironmentConfigIns().LogLevel)
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	//fmt.Printf("'%s'",level)
	if level == "debug" {
		atomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if level == "info" {
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else if level == "warn" {
		atomicLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	} else if level == "error" {
		atomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else if level == "panic" {
		atomicLevel = zap.NewAtomicLevelAt(zap.PanicLevel)
	} else if level == "fatal" {
		atomicLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	} else if level == "dpanic" {
		atomicLevel = zap.NewAtomicLevelAt(zap.DPanicLevel)
	} else {
		atomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	return zap.Config{
		Level:       atomicLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    newZapEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func init() {
	loggerOnce.Do(func() {
		loggerIns = &Factory{}
		loggerIns.setup()
		err := loggerIns.buildZapLogger()
		if err != nil {
			panic(err)
		}
	})
}