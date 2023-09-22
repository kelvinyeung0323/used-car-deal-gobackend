package config

import (
	"log"

	"github.com/spf13/viper"
)

var MySQLConfig = struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}{}

var AppConfig = struct {
	Port string `yaml:"port"`
}{}

var LogConfig = struct {
	ConsoleOnly bool   `json:"consoleOnly"`
	Level       string `json:"level"`
	Filename    string `json:"filename"`
	MaxSize     int    `json:"maxsize"`
	MaxAge      int    `json:"max_age"`
	MaxBackups  int    `json:"max_backups"`
}{}

func InitConfig() {
	//TODO:监控配置文件变化
	loadAppConfig()
}

func loadAppConfig() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf/")

	if err := viper.ReadInConfig(); err != nil {
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// 	// Config file not found; ignore error if desired
		// } else {
		// 	// Config file was found but another error was produced
		// }
		log.Printf("error:读取文件失败,%v", err)
		panic("读取配置文件失败.")
	}

	//解释日志配置
	logger := viper.Sub("logger")

	if err := logger.Unmarshal(&LogConfig); err != nil {
		log.Printf("error:解释配置文件失败,%v", err)
	}

	//解释Mysql配置
	mysql := viper.Sub("datasource.mysql")

	if err := mysql.Unmarshal(&MySQLConfig); err != nil {
		log.Printf("error:解释配置文件失败,%v", err)
	}

	//解释APP配置
	app := viper.Sub("app")

	if err := app.Unmarshal(&AppConfig); err != nil {
		log.Printf("error:解释配置文件失败,%v", err)
	}

}
