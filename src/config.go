package main

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"reflect"
)

var (
	defConfigPath      = "."
	defConfigName      = "ftpdts.ini"
	defConfigType      = "ini"
	defConfigEnvPrefix = "FTPDTS"
)

type Config struct {
	HTTP struct {
		Port           uint   `default:"2001"`
		Host           string `default:"127.0.0.1"`
		MaxRequestBody int64  `default:"1024"`
	}

	FTP struct {
		Port         uint   `default:"2000"`
		Host         string `default:"127.0.0.1"`
		PassivePorts string `default:"32000-32010"`
		DebugMode    bool   `default:"false"`
		PublicIP     string
	}

	Templates struct {
		Path string `default:"./tmpl"`
	}

	Data struct {
		Path string `default:"./data"`
	}

	Logs struct {
		FTP             string `default:"logs/ftp.log"`
		FTPNoConsole    bool   `default:"false"`
		HTTP            string `default:"logs/http.log"`
		HTTPNoConsole   bool   `default:"false"`
		Ftpdts          string `default:"logs/ftpdts.log"`
		FtpdtsNoConsole bool   `default:"false"`
	}

	Cache struct {
		DataTTL uint `default:"86400"`
	}

	UID struct {
		Chars           string `default:"1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"`
		Format          string `default:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`
		ValidatorRegexp string `default:"[0-9a-zA-Z]{32}"`
	}
}

func viperConfig(cPath string, cName string, cType string, envPrefix string, config interface{}) (v *viper.Viper, err error) {
	v = viper.New()
	v.AddConfigPath(cPath)
	v.SetConfigType(cType)
	v.AutomaticEnv()
	v.SetEnvPrefix(envPrefix)
	v.SetConfigName(cName)
	//fill default config values
	if reflect.TypeOf(config).Kind() == reflect.Ptr && reflect.ValueOf(config).Elem().Type().Kind() == reflect.Struct {
		if err = defaults.Set(config); err != nil {
			return
		}
	}
	return
}

func LoadConfig() (c *Config, err error) {

	c = new(Config)
	viperConfig, err := viperConfig(defConfigPath, defConfigName, defConfigType, defConfigEnvPrefix, c)
	if err != nil {
		return
	}

	// Find and read the config file
	if err := viperConfig.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("wrong config file: %v", err)
	}

	if err := viperConfig.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}

	return
}
