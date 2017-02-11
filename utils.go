package config

import "os"

func getEnv(name string) string {
	e := os.Getenv(name)
	if e == "" {
		return ""
	}
	return e
}
