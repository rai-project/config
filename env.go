package config

import (
	"os"
	"path/filepath"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/mitchellh/go-homedir"

	"github.com/rai-project/godotenv"
)

func initEnv() {
	if IsVerbose || IsDebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// load exePath/.rai.profile or exePath/rai.profile where exePath is the path of the rai executable
	exeDir := filepath.Dir(os.Args[0])
	if dir, err := filepath.Abs(exeDir); err == nil {
		currentEnvFile := filepath.Join(dir, "."+DefaultAppName+".profile")
		if com.IsFile(currentEnvFile) {
			log.WithField("env_file", currentEnvFile).Debug("reading environment from current directory env file")
			godotenv.Overload(currentEnvFile)
		}
		currentEnvFile = filepath.Join(dir, DefaultAppName+".profile")
		if com.IsFile(currentEnvFile) {
			log.WithField("env_file", currentEnvFile).Debug("reading environment from current directory env file")
			godotenv.Overload(currentEnvFile)
		}
	}

	// load ~/.rai.env
	homeEnvFile, err := homedir.Expand("~/." + DefaultAppName + ".env")
	if err == nil && com.IsFile(homeEnvFile) {
		log.WithField("env_file", homeEnvFile).Debug("reading environment from home directory env file")
		godotenv.Overload(homeEnvFile)
	}
	// load ~/.rai.profile
	homeEnvFile, err = homedir.Expand("~/." + DefaultAppName + ".profile")
	if err == nil && com.IsFile(homeEnvFile) {
		log.WithField("env_file", homeEnvFile).Debug("reading environment from home directory profile file")
		godotenv.Overload(homeEnvFile)
	}

	srcpath, err := sourcepath.AbsoluteDir()
	baseDir := filepath.Dir(filepath.Dir(srcpath))
	if err == nil {
		envFile := filepath.Join(baseDir, ".env")
		if com.IsFile(envFile) {
			log.WithField("env_file", envFile).Debug("reading environment from local env file")
			godotenv.Overload(envFile)
		}

		privateEnvFile := filepath.Join(baseDir, ".env.private")
		if com.IsFile(privateEnvFile) {
			log.WithField("env_file", privateEnvFile).Debug("reading environment from private env file")
			godotenv.Overload(privateEnvFile)
		}
	}
}
