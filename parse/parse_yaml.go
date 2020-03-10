package parse

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Yaml 解析yaml文件
type Yaml struct {
	filepath  string
	filename  string
	viperSeed *viper.Viper
}

// initViper 初始化viper
func initViper(filepath string, filename string) *viper.Viper {
	// Init viper path and filename/type

	viperSeed := viper.New()
	filenameSplit := strings.Split(filename, ".")
	if len(filenameSplit) != 2 {
		viperSeed.SetConfigType("yaml")
	} else {
		viperSeed.SetConfigType(filenameSplit[1])
		viperSeed.SetConfigName(filenameSplit[0])
	}
	viperSeed.AddConfigPath(filepath)

	if err := viperSeed.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("No such config file")
		} else {
			log.Fatal("Read config error")
		}
	}
	log.Println("Init viper complete")
	return viperSeed
}

// NewYaml NewYaml
func NewYaml(filepath string, filename string) *Yaml {
	viperSeed := initViper(filepath, filename)
	return &Yaml{
		filepath:  filepath,
		filename:  filename,
		viperSeed: viperSeed,
	}
}

// Parse 解析yaml
func (yaml *Yaml) Parse() map[string]interface{} {
	return yaml.viperSeed.AllSettings()
}
