package config

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
)

const (
	DefaultAppName        = "rai"
	DefaultAppFullName    = "rai project"
	DefaultAppDescription = ""
	DefaultLicense        = "NCSA"
	DefaultURL            = "rai-project.com"
)

// APP holds common application fields credentials and keys.
type appConfig struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	License     string `json:"license"`
	URL         string `json:"url"`
	IsDebug     bool   `json:"debug"`
	IsVerbose   bool   `json:"verbose"`
}

func (appConfig) ConfigName() string {
	return "App"
}

func (appConfig) setDefaults() {
	if viper.Get("app") == nil {
		viper.SetDefault("app", map[string]interface{}{})
	}

	viper.SetDefault("app.name", DefaultAppName)
	viper.SetDefault("app.full_name", DefaultAppFullName)
	viper.SetDefault("app.description", DefaultAppDescription)
	viper.SetDefault("app.license", DefaultLicense)
	viper.SetDefault("app.url", DefaultURL)
	viper.SetDefault("app.debug", false)
	viper.SetDefault("app.verbose", false)
}

func (a *appConfig) Read() {
	IsDebug = viper.GetBool("app.debug")
	IsVerbose = viper.GetBool("app.verbose")
	*App = appConfig{
		Name:        viper.GetString("app.name"),
		FullName:    viper.GetString("app.full_name"),
		Description: viper.GetString("app.description"),
		License:     viper.GetString("app.license"),
		URL:         viper.GetString("app.url"),
		IsDebug:     IsDebug,
		IsVerbose:   IsVerbose,
	}
}

func (c appConfig) String() string {
	return pp.Sprintln(c)
}

func (c appConfig) Debug() {
	log.Debug("App Config = ", c)
}
