package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"bytes"

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

func setViperConfig(opts *Options) {
	if opts.ConfigString != "" {
		return
	}

	defer viper.AutomaticEnv() // read in environment variables that match
	defer viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if com.IsFile(opts.ConfigFileAbsolutePath) {
		log.Debug("Found ", opts.ConfigFileAbsolutePath, " already set. Using ", opts.ConfigFileAbsolutePath, " as the config file.")
		viper.SetConfigFile(opts.ConfigFileAbsolutePath)
		return
	}
	if val, ok := os.LookupEnv(opts.ConfigEnvironName); ok {
		pth, _ := homedir.Expand(val)
		log.Debug("Found ", opts.ConfigEnvironName, " in env. Using ", val, " as config file name")
		if com.IsFile(pth) {
			viper.SetConfigFile(pth)
			return
		}
		dir, file := path.Split(pth)
		ext := path.Ext(file)
		file = strings.TrimSuffix(file, ext)
		viper.SetConfigName(file)
		viper.AddConfigPath(dir)
		return
	}
	if pth, err := homedir.Expand("~/." + opts.AppName + "_config.yaml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ~/." + opts.AppName + "_config.yaml as config file.")
		viper.SetConfigFile(pth)
		return
	}
	if pth, err := filepath.Abs("../." + opts.AppName + "_config.yaml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ../." + opts.AppName + "_config.yaml as config file.")
		viper.SetConfigFile(pth)
		return
	}

	defer func() {
		for _, pth := range opts.ConfigSearchPaths {
			pth, err := homedir.Expand(pth)
			if err != nil {
				continue
			}
			viper.AddConfigPath(pth)
		}
		viper.SetConfigType(opts.ConfigFileType)
	}()

	log.Info("No fixed configuration file found, searching for a config file with name=", opts.ConfigFileBaseName)
	viper.SetConfigName(opts.ConfigFileBaseName)
}

func load(opts *Options) {
	initEnv(opts)
	if opts.ConfigString == "" {
		setViperConfig(opts)

		// read configuration
		err := viper.ReadInConfig()
		if err != nil {
			log.WithError(err).
				WithField("config_file", viper.ConfigFileUsed()).
				Error("Cannot read in configuration file ")
		}
	} else {
		reader := bytes.NewBufferString(opts.ConfigString)
		err := viper.ReadConfig(reader)
		if err != nil {
			log.WithError(err).
				Error("Cannot read in configuration string ")
		}
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
