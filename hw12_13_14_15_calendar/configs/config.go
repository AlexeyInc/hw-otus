package calendarconfig

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	Storage    StorageConf
	HTTPServer ServerConf
	GRPCServer ServerConf
	AMQP       AMQPConf
	Scheduler  SchedulerConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Driver string
	Source string
}

type ServerConf struct {
	Host    string
	Port    string
	Network string
}

type AMQPConf struct {
	Source string
}

type SchedulerConf struct {
	CheckNotificationFreqSeconds  int
	CheckExpiredEventsFreqSeconds int
}

func NewConfig(fullPath string) (config Config, err error) {
	pathToFile, fileName, fileType := getFileInfo(fullPath)
	viper.AddConfigPath(pathToFile)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func getFileInfo(filePath string) (path, fileName, fileType string) {
	pathChanks := strings.Split(filePath, "/")

	sb := strings.Builder{}
	for i := 0; i < len(pathChanks)-1; i++ {
		sb.WriteString(pathChanks[i])
		sb.WriteString("/")
	}
	path = sb.String()
	fileName = strings.Split(pathChanks[len(pathChanks)-1], ".")[0]
	fileType = strings.Split(pathChanks[len(pathChanks)-1], ".")[1]
	return
}
