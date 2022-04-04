package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     loggerConf
	Storage    storageConf
	HTTPServer serverConf
	GRPCServer serverConf
}

type loggerConf struct {
	Level string `toml:"level"`
}

type storageConf struct {
	Driver string `toml:"driver"`
	Source string `toml:"source"`
}

type serverConf struct {
	Host    string `toml:"host"`
	Port    string `toml:"port"`
	Network string `toml:"network"`
}

func NewConfig(fullPath string) (config Config, err error) {
	pathToFile, fileName, fileType := getFileInfo(fullPath)
	viper.AddConfigPath(pathToFile)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)

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
