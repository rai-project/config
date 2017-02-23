package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type ConfigInterface interface {
	ConfigName() string
	SetDefaults()
	Read()
	String() string
	Debug()
}

var (
	readMutex sync.Mutex
	IsVerbose = false
	IsDebug   = false
	log       = logrus.WithField("pkg", "config")
)

var (
	ConfigPaths       = []string{"$HOME", "..", "../..", "."}
	ConfigEnvironName = "RAI_CONFIG_FILE"
	ConfigFileName    = ".rai_config"
	ConfigFileType    = "yaml"
)

func setViperConfig() {
	defer viper.AutomaticEnv() // read in environment variables that match
	defer viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	defer func() {
		for _, pth := range ConfigPaths {
			if pth[0] != '$' {
				viper.AddConfigPath(pth)
				continue
			}
			if val, ok := os.LookupEnv(pth); ok {
				viper.AddConfigPath(val)
				continue
			}
			viper.AddConfigPath(pth)
		}

		viper.SetConfigType(ConfigFileType)
	}()

	ConfigFileName = "." + appName + "_config"
	ConfigEnvironName = strings.ToUpper(appName) + "_CONFIG_FILE"

	if com.IsFile(ConfigFileName) {
		log.Debug("Found ", ConfigFileName, " being set. Using ", ConfigFileName, " as the config file.")
		viper.SetConfigFile(ConfigFileName)
		return
	}
	if val, ok := os.LookupEnv(ConfigEnvironName); ok {
		pth, _ := homedir.Expand(val)
		log.Debug("Found ", ConfigEnvironName, " in env. Using ", val, " as config file name")
		if com.IsFile(pth) {
			viper.SetConfigFile(ConfigFileName)
			return
		}
		dir, file := path.Split(pth)
		ext := path.Ext(file)
		file = strings.TrimSuffix(file, ext)
		viper.SetConfigName(file)
		viper.AddConfigPath(dir)
		return
	}
	if pth, err := homedir.Expand("~/." + appName + "_config.yaml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ~/." + appName + "_config.yaml as config file.")
		viper.SetConfigFile(pth)
		return
	}
	if pth, err := filepath.Abs("../." + appName + "_config.yaml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ../." + appName + "_config.yaml as config file.")
		viper.SetConfigFile(pth)
		return
	}

	log.Info("No fixed configuration file found, searching for a config file with name=", ConfigFileName)
	viper.SetConfigName(ConfigFileName)
}

func load() {

	readMutex.Lock()
	defer readMutex.Unlock()

	initEnv()
	setViperConfig()

	// read configuration
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).
			WithField("config_file", viper.ConfigFileUsed()).
			Error("Cannot read in configuration file ")
		return
	}

	for _, r := range registry {
		r.SetDefaults()
	}
	for _, r := range registry {
		r.Read()
	}
}

func Debug() {
	log.Debug("Config = ")
	for _, r := range registry {
		r.Debug()
	}
}
