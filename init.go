package config

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/k0kubun/pp"
)

var (
	log             *logrus.Entry
	once            sync.Once
	onInitFunctions struct {
		funcs []func()   `json:"funcs"`
		mutex sync.Mutex `json:"mutex"`
	}
	afterInitFunctions struct {
		funcs []func()   `json:"funcs"`
		mutex sync.Mutex `json:"mutex"`
	}
)

func OnInit(f func()) {
	onInitFunctions.mutex.Lock()
	defer onInitFunctions.mutex.Unlock()
	onInitFunctions.funcs = append(onInitFunctions.funcs, f)
}

func AfterInit(f func()) {
	afterInitFunctions.mutex.Lock()
	defer afterInitFunctions.mutex.Unlock()
	afterInitFunctions.funcs = append(afterInitFunctions.funcs, f)
}

func Init(opts ...Option) {
	once.Do(func() {
		options := NewOptions()

		for _, o := range opts {
			o(options)
		}

		if options.AppSecret != "" {
			defer func() {
				App.Secret = options.AppSecret
			}()
		}

		if options.AppName != "" {
			defer func() {
				App.Name = options.AppName
			}()
		}

		log = logrus.WithField("pkg", "config")
		if options.IsDebug || options.IsVerbose {
			pp.WithLineInfo = true
			log.Level = logrus.DebugLevel
		}

		load(options)

		if initFunsLength := len(onInitFunctions.funcs); initFunsLength > 0 {
			var wg sync.WaitGroup
			wg.Add(initFunsLength)
			for ii := range onInitFunctions.funcs {
				f := onInitFunctions.funcs[ii]
				go func() {
					defer wg.Done()
					f()
				}()
			}
			wg.Wait()
		}

		if initFunsLength := len(afterInitFunctions.funcs); initFunsLength > 0 {
			var wg sync.WaitGroup
			wg.Add(initFunsLength)
			for ii := range afterInitFunctions.funcs {
				f := afterInitFunctions.funcs[ii]
				go func() {
					defer wg.Done()
					f()
				}()
			}
			wg.Wait()
		}

		if IsVerbose {
			fmt.Println("Finished setting configuration...")
		}
	})
}

func init() {
	log = logrus.WithField("pkg", "config")

	opts := NewOptions()
	if opts.IsVerbose || opts.IsDebug {
		log.Level = logrus.DebugLevel
	}

	initEnv(opts)

	isVerbose, isDebug := modeInfo()
	if isVerbose || isDebug {
		log.Level = logrus.DebugLevel
	}
}
