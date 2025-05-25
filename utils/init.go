package utils

import (
	"fmt"

	"github.com/jiny3/gopkg/configx"
	"github.com/jiny3/gopkg/envx"
	"github.com/jiny3/gopkg/filex"
	"github.com/jiny3/gopkg/logx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Init = func() {
	confPath, err := envx.Home()
	if err != nil {
		logrus.WithError(err).Warn("get home path failed, use default config path")
		err = configx.Load("config/config.yaml")
		if err != nil {
			logx.Init(logrus.DebugLevel)
			logrus.WithError(err).Warn("config load failed")
		}
	} else {
		confPath = fmt.Sprintf("%s/.config/xAI/config.yaml", confPath)
		_, err = filex.FileCreate(confPath)
		if err != nil {
			logrus.WithError(err).Warn("create at .config/ failed, use default config path")
			err = configx.Load("config/config.yaml")
			if err != nil {
				logx.Init(logrus.DebugLevel)
				logrus.WithError(err).Warn("config load failed")
			}
		} else {
			err = configx.Load(confPath)
			if err != nil {
				logx.Init(logrus.DebugLevel)
				logrus.WithError(err).Warn("config load failed")

			}
		}
	}

	logPath, level := viper.GetString("log.path"), viper.GetString("log.level")
	_level, err := logrus.ParseLevel(level)
	if err != nil {
		_level = logrus.InfoLevel
	}
	if logPath == "" {
		logx.Init(_level)
	} else {
		logx.Init(_level, logPath)
	}
}
