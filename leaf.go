package leaf

import (
	"github.com/abel/leaf/cluster"
	"github.com/abel/leaf/conf"
	"github.com/abel/leaf/console"
	"github.com/abel/leaf/log"
	"github.com/abel/leaf/module"
	"os"
	"os/signal"
)

var (
	stopChan = make(chan int)
)

func Run(mods ...module.Module) {
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	log.Release("Leaf %v starting up", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case sig := <-c:
		log.Release("Leaf closing down (signal: %v)", sig)
	case reason := <-stopChan:
		log.Release("Leaf closing down (reason: %v)", reason)
	}
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}

func Stop(reason int) {
	stopChan <- reason
}
