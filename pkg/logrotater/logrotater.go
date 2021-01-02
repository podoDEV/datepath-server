package logrotater

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

const (
	defaultLogMaxSize    = 300
	defaultLogMaxBackups = 30
	defaultLogMaxAge     = 30
	defaultLogCompress   = false
)

type Config struct {
	MaxSize    int  `toml:"log_max_size"`
	MaxBackups int  `toml:"log_max_backups"`
	MaxAge     int  `toml:"log_max_age"`
	Compress   bool `toml:"log_compress"`
}

func NewConfig() Config {
	return Config{
		MaxSize:    defaultLogMaxSize,
		MaxBackups: defaultLogMaxBackups,
		MaxAge:     defaultLogMaxAge,
		Compress:   defaultLogCompress,
	}
}

type LogRotater struct {
	Logger lumberjack.Logger
}

func New(path string, conf *Config) LogRotater {
	return LogRotater{
		Logger: lumberjack.Logger{
			Filename:   path,
			MaxSize:    conf.MaxSize, // megabytes
			MaxBackups: conf.MaxBackups,
			MaxAge:     conf.MaxAge,   // days
			Compress:   conf.Compress, // disabled by default
		},
	}
}

func (l *LogRotater) Run() {
	for {
		t := <-time.Tick(time.Hour)
		if t.Hour() == 0 {
			l.Logger.Rotate()
		}
	}
}
