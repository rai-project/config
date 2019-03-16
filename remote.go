//+build remote_config

package config

import (
	"strings"

	_ "github.com/spf13/viper/remote"
)

var validRemotePrefixes = []string{
	"etcd://",
	"consul://",
}

// IsValidRemotePrefix ...
func IsValidRemotePrefix(s string) bool {
	for _, p := range validRemotePrefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
