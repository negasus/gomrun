package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fatih/color"

	"github.com/negasus/gomrun/internal/config"
	"github.com/negasus/gomrun/internal/exec"
)

var version = "v0.1.0"

var (
	configFile string
	colors     = []color.Attribute{
		color.FgGreen,
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
	}
)

func main() {
	fmt.Printf("%s %s\n", color.CyanString("gomrun"), color.GreenString(version))

	flag.StringVar(&configFile, "config", ".gomrun.yml", "config file")
	flag.Parse()

	cfgData, errReadConfig := os.ReadFile(configFile)
	if errReadConfig != nil {
		color.Red("error read config file %q, %v\n", configFile, errReadConfig)
		os.Exit(1)
	}

	cfg, errConfig := config.Load(cfgData)
	if errConfig != nil {
		color.Red("error load config %q, %v\n", configFile, errConfig)
		os.Exit(1)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	services := map[string]struct{}{}
	for _, r := range flag.Args() {
		services[r] = struct{}{}
	}

	var nameMaxLen int
	for name := range cfg.Services {
		if len(name) > nameMaxLen {
			nameMaxLen = len(name)
		}
	}

	wg := sync.WaitGroup{}

	var i int
	for name := range cfg.Services {
		if len(services) > 0 {
			if _, ok := services[name]; !ok {
				continue
			}
		}
		wg.Add(1)
		go exec.Exec(&wg, cfg, name, nameMaxLen, colors[i%len(colors)])
		i++
	}

	wg.Wait()
}
