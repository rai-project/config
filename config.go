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
	setDefaults()
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
	ConfigFileName    = "rai_config"
	ConfigFileType    = "yaml"
)

func loadViper() {
	defer viper.AutomaticEnv() // read in environment variables that match
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
	if pth, err := homedir.Expand("~/.rai.yaml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ~/.rai.yaml as config file.")
		home, _ := homedir.Dir()
		viper.SetConfigFile(pth)
		viper.SetConfigName(".rai")
		viper.AddConfigPath(home)
		return
	}
	if pth, err := filepath.Abs("../rai_config.yaml"); err == nil && com.IsFile(pth) {
		log.Debug("Using \"" + pth + "\" as config file.")
		viper.SetConfigFile(pth)
		viper.SetConfigName("rai_config")
		viper.AddConfigPath(filepath.Dir(pth))
		return
	}

	log.Info("No fixed configuration file found, searching for a config file with name=", ConfigFileName)
	viper.SetConfigName(ConfigFileName)
}

func load() {

	readMutex.Lock()
	defer readMutex.Unlock()

	initEnv()
	loadViper()

	// read configuration
	err := viper.ReadInConfig()
	if err == nil && viper.IsSet("config") {
		viper.SetConfigFile(viper.GetString("config"))
		if com.IsFile(viper.GetString("config")) {
			// If a config file is found, read it in.
			viper.AutomaticEnv() // read in environment variables that match
			err = viper.ReadInConfig()
		}
	}
	if err != nil {
		log.WithError(err).Panic("Cannot read in configuration file")
	}

	if err != nil {
		viper.Debug()
		log.Panic("Cannot read configuration file. ")
	}
	if IsVerbose {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
	for _, r := range registry {
		r.setDefaults()
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
