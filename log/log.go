package log

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *Logger

	G = Get

	defaultConfig = Config{
		Level:         "info",
		Format:        "console",
		HideFilename:  true,
		HideTimestamp: true,
	}
)

type (
	// Config is the logger configuration
	Config struct {
		Level         string
		Format        string
		HideFilename  bool
		HideTimestamp bool
	}

	Logger struct {
		*zap.SugaredLogger
	}
)

func init() {
	var err error

	logger, err = Build(&defaultConfig)
	if err != nil {
		panic(err)
	}
}

// Build builds a new logger from a config. If conf is nil, uses default configuration.
func Build(conf *Config) (*Logger, error) {
	var (
		err error
		l   *zap.Logger
	)

	if conf == nil {
		l, err = zap.NewDevelopment()
	} else {
		var zc *zap.Config
		zc, err = conf.parseZap()
		if err != nil {
			return nil, err
		}

		l, err = zc.Build()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build logger from config: %w", err)
	}

	logger = &Logger{l.Sugar()}

	return logger, nil
}

// AddFlags binds the logger configuration flags to the provided flag set and returns a
// config object that can return a configured logger by calling the Build() function.
func AddFlags(fs *pflag.FlagSet) *Config {
	conf := &Config{}

	fs.StringVar(&conf.Level, "log-level", "info", "one of: info|debug|warn")
	fs.StringVar(&conf.Format, "log-format", "console", "one of: json|console")

	return conf
}

// Get returns the default configured logger. You may also use its shorter alias: @G()
func Get() *Logger {
	if logger == nil {
		panic("logger not yet configured, need to call Build() first!")
	}

	return logger
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{l.SugaredLogger.Named(name)}
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{l.SugaredLogger.With(args...)}
}

// NewMiddleware returns an http middleware that logs requests
func (l *Logger) NewMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		method := r.Method
		path := r.URL.Path
		rwSpy := &responseWriterSpy{rw: w}

		next(rwSpy, r)

		duration := time.Since(startTime)

		l.Infow("finished http request",
			"method", method,
			"path", path,
			"code", rwSpy.statusCode,
			"duration", duration,
		)
	}
}

// implement io.Writer
func (l *Logger) Write(data []byte) (int, error) {
	l.Error(string(data))
	return len(data), nil
}

type responseWriterSpy struct {
	statusCode int
	rw         http.ResponseWriter
}

func (rw *responseWriterSpy) Header() http.Header {
	return rw.rw.Header()
}

func (rw *responseWriterSpy) Write(data []byte) (int, error) {
	return rw.rw.Write(data)
}

func (rw *responseWriterSpy) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.rw.WriteHeader(statusCode)
}

func (c *Config) parseZap() (*zap.Config, error) {
	ret := zap.NewDevelopmentConfig()
	ret.DisableStacktrace = true

	switch c.Level {
	case "info":
		ret.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "debug":
		ret.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		ret.DisableStacktrace = false
	case "warn":
		ret.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		return nil, fmt.Errorf("unknown log level: %s", c.Level)
	}

	switch c.Format {
	case "console":
		ret.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		ret.EncoderConfig.ConsoleSeparator = " "

		if c.HideTimestamp {
			ret.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {}
		} else {
			ret.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("06-02-01T15:04:05.000")
		}

		if c.HideFilename {
			ret.EncoderConfig.EncodeCaller = func(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {}
		} else {
			ret.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, en zapcore.PrimitiveArrayEncoder) {
				en.AppendString(caller.TrimmedPath() + "]")
			}
		}

	case "json":
		ret.Encoding = c.Format
	default:
		return nil, fmt.Errorf("unknown log format: %s", c.Format)
	}

	return &ret, nil
}
