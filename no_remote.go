//+build !remote_config

package config

var validRemotePrefixes = []string{}

func IsValidRemotePrefix(s string) bool {
	return false
}
