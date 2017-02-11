package config

import "sync"

var mutex sync.Mutex

func Register(mod ConfigInterface) {
	mutex.Lock()
	defer mutex.Unlock()
	modules = append(modules, mod)
}
