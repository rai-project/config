package config

import (
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	colorable "github.com/mattn/go-colorable"
	"github.com/rai-project/vipertags"
	"github.com/spf13/viper"
)

const (
	DefaultAppDescription = ""
)

// APP holds common application fields credentials and keys.
type appConfig struct {
	Name        string `json:"name" config:"app.name" default:"default"`
	FullName    string `json:"full_name" config:"app.full_name" default:"rai project"`
	Description string `json:"description" config:"app.description"`
	License     string `json:"license" config:"app.license" default:"NCSA"`
	URL         string `json:"url" config:"app.url" default:"rai-project.com"`
	Secret      string `json:"-" config:"app.secret" default:"default"`
	Color       bool   `json:"color" config:"app.color" env:"COLOR"`
	IsDebug     bool   `json:"debug" config:"app.debug" env:"DEBUG"`
	IsVerbose   bool   `json:"verbose" config:"app.verbose" env:"VERBOSE"`
}

var (
	App              = &appConfig{}
	DefaultAppName   = "rai"
	DefaultAppSecret = "-secret-"
	DefaultAppColor  = !color.NoColor
	IsDebug          bool
	IsVerbose        bool
)

func (appConfig) ConfigName() string {
	return "App"
}

func (a *appConfig) SetDefaults() {
	viper.Set("app.color", DefaultAppColor)
}

func (a *appConfig) Read() {
	IsDebug = viper.GetBool("app.debug") || IsDebug
	IsVerbose = viper.GetBool("app.verbose") || IsVerbose
	vipertags.Fill(a)
	if a.Name == "" || a.Name == "default" {
		a.Name = DefaultAppName
	}
	if !viper.IsSet("app.color") {
		a.Color = DefaultAppColor
		viper.Set("app.color", a.Color)
	}
	if a.Secret == "" || a.Secret == "default" {
		a.Secret = DefaultAppSecret
	}
	if a.Color == false {
		pp.SetDefaultOutput(colorable.NewNonColorable(pp.GetDefaultOutput()))
	}
	if a.IsDebug || a.IsVerbose {
		pp.WithLineInfo = true
	}
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
