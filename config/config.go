package config

import (
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"
)

var config C

type C struct {
	Env      string
	Database struct {
		Driver     string
		Connection string
		Logger     struct {
			Kafka echomiddleware.KafkaConfig
		}
	}
	BehaviorLog struct {
		Kafka echomiddleware.KafkaConfig
	}
	Services struct {
		ColleagueApi string
	}
	ServiceName string
	HttpPort    string
	AppEnv      string
}

func Init(appEnv string, options ...func(*C)) C {
	config.AppEnv = appEnv
	if err := configutil.Read(appEnv, &config); err != nil {
		logrus.WithError(err).Warn("Fail to load config file")
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

func Config() C {
	return config
}
