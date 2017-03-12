package config

import (
	"strings"

	"github.com/fatih/color"
)

type Options struct {
	AppName                string
	AppSecret              string
	ConfigSearchPaths      []string
	ConfigEnvironName      string
	ConfigFileBaseName     string
	ConfigFileType         string
	ConfigFileAbsolutePath string
	ConfigString           string
	IsVerbose              bool
	IsDebug                bool
}

type Option func(*Options)

func NewOptions() *Options {
	isVerbose, isDebug := modeInfo()
	return &Options{
		AppName:            DefaultAppName,
		ConfigSearchPaths:  []string{"$HOME", "..", "../..", "."},
		ConfigEnvironName:  strings.ToUpper(DefaultAppName) + "_CONFIG_FILE",
		ConfigFileBaseName: "." + strings.ToLower(DefaultAppName) + "_config",
		ConfigFileType:     "yaml",
		IsDebug:            isDebug,
		IsVerbose:          isVerbose,
	}
}

func AppName(s string) Option {
	return func(opts *Options) {
		DefaultAppName = s
		opts.AppName = s
		opts.ConfigFileBaseName = "." + strings.ToLower(DefaultAppName) + "_config"
	}
}

func AppSecret(s string) Option {
	return func(opts *Options) {
		DefaultAppSecret = s
		opts.AppSecret = s
	}
}

func ConfigSearchPaths(s []string) Option {
	return func(opts *Options) {
		opts.ConfigSearchPaths = s
	}
}

func ConfigEnvironName(s string) Option {
	return func(opts *Options) {
		opts.ConfigEnvironName = s
	}
}

func ConfigFileBaseName(s string) Option {
	return func(opts *Options) {
		opts.ConfigFileBaseName = s
	}
}

func ConfigFileType(s string) Option {
	return func(opts *Options) {
		opts.ConfigFileType = s
	}
}

func ConfigFileAbsolutePath(s string) Option {
	return func(opts *Options) {
		opts.ConfigFileAbsolutePath = s
	}
}

func ConfigString(s string) Option {
	return func(opts *Options) {
		opts.ConfigString = s
	}
}

func VerboseMode(b bool) Option {
	return func(opts *Options) {
		opts.IsVerbose = b
		IsVerbose = b
	}
}

func DebugMode(b bool) Option {
	return func(opts *Options) {
		opts.IsDebug = b
		IsDebug = b
	}
}

func ColorMode(b bool) Option {
	return func(opts *Options) {
		DefaultAppColor = b
		color.NoColor = !b
	}
}
