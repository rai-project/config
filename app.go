package config

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/vipertags"
	"github.com/spf13/viper"
)

const (
	DefaultAppDescription = ""
)

// APP holds common application fields credentials and keys.
type appConfig struct {
	Name        string `json:"name" config:"app.name" default:"rai"`
	FullName    string `json:"full_name" config:"app.full_name" default:"rai project"`
	Description string `json:"description" config:"app.description"`
	License     string `json:"license" config:"app.license" default:"NCSA"`
	URL         string `json:"url" config:"app.url" default:"rai-project.com"`
	IsDebug     bool   `json:"debug" config:"app.debug" env:"DEBUG"`
	IsVerbose   bool   `json:"verbose" config:"app.verbose" env:"VERBOSE"`
}

var (
	App = &appConfig{}
)

func (appConfig) ConfigName() string {
	return "App"
}

func (appConfig) SetDefaults() {
}

func (a *appConfig) Read() {
	IsDebug = viper.GetBool("app.debug")
	IsVerbose = viper.GetBool("app.verbose")
	vipertags.Fill(a)
}

func (a appConfig) String() string {
	return pp.Sprintln(a)
}

func (a appConfig) Debug() {
	log.Debug("App Config = ", a)
}

func init() {
	Register(App)
}
