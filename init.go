package config

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/k0kubun/pp"
)

var (
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

func Init() {
	once.Do(func() {

		log = logrus.WithField("pkg", "config")

		load()

		if IsDebug {
			pp.WithLineInfo = true
		}

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
		// if Mode.IsVerbose {
		// 	fmt.Println("Finished running configuration...")
		// }
	})
}

func init() {
	initEnv()
}
