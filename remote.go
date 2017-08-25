package config

import "strings"

var validRemotePrefixes = []string{
	"etcd://",
	"consul://",
}

func IsValidRemotePrefix(s string) bool {
	for _, p := range validRemotePrefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
