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
	IsServer  = false
	IsClient  = false
	log       = logrus.WithField("pkg", "config")

	App     = new(appConfig)
	modules = []ConfigInterface{
		App,
	}
)

var (
	ConfigPaths       = []string{"$HOME", "..", "../..", "."}
	ConfigEnvironName = "RAI_CONFIG_FILE"
	ConfigFileName    = "rai_config"
	ConfigFileType    = "yaml"
)

func Read() {

	readMutex.Lock()
	defer readMutex.Unlock()

	initEnv()

	if com.IsFile(ConfigFileName) {
		dir, file := path.Split(ConfigFileName)
		ext := path.Ext(file)
		file = strings.TrimSuffix(file, ext)
		viper.SetConfigName(file)
		viper.AddConfigPath(dir)
	} else if val, ok := os.LookupEnv(ConfigEnvironName); ok {
		log.Info("Found ", ConfigEnvironName, " in env. Using ", val, " as config file name")
		pth, _ := homedir.Expand(val)
		dir, file := path.Split(pth)
		ext := path.Ext(file)
		file = strings.TrimSuffix(file, ext)
		viper.SetConfigName(file)
		viper.AddConfigPath(dir)
	} else if pth, err := homedir.Expand("~/.rai.yaml"); err == nil && com.IsFile(pth) {
		log.Info("Using ~/.rai.yaml as config file.")
		home, _ := homedir.Dir()
		viper.SetConfigName(".rai")
		viper.AddConfigPath(home)
	} else if pth, err := filepath.Abs("../rai_config.yaml"); err == nil && com.IsFile(pth) {
		log.Info("Using \"" + pth + "\" as config file.")
		viper.SetConfigName("rai_config")
		viper.AddConfigPath(filepath.Dir(pth))
	} else {
		log.Info("No fixed configuration file found, searching for a config file with name=", ConfigFileName)
		viper.SetConfigName(ConfigFileName)
	}

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
	// If a config file is found, read it in.
	viper.AutomaticEnv() // read in environment variables that match

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
	for _, mod := range modules {
		mod.setDefaults()
	}
	for _, mod := range modules {
		mod.Read()
	}
}

func Debug() {
	log.Debug("Config = ")
	for _, mod := range modules {
		mod.Debug()
	}
}
