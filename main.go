package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/gdamore/tcell"
	"github.com/golang/glog"
	"github.com/tupyy/tview"
	"github.com/tupyy/lazylogger/internal/conf"
	"github.com/tupyy/lazylogger/internal/gui"
	"github.com/tupyy/lazylogger/internal/log"
)

// build flags
var (
	Version string

	Build string

	BuildDate string

	configurationFile string

	// app
	app = tview.NewApplication()

	loggerManager *log.LoggerManager
)

func main() {

	var version = flag.Bool("version", false, "Show version")

	// Read configuration
	flag.StringVar(&configurationFile, "config", "nodata", "JSON configuration file")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *version {
		fmt.Printf("LazyLogger:\n %-8s: %-10s\n %-8s: %-10s\n %-8s: %-10s\n",
			"Version", Version,
			"Build", Build,
			"Date", BuildDate)

		os.Exit(0)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			glog.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	config := conf.ReadConfigurationFile(configurationFile)
	glog.Infof("Configuration has %d.", len(config.LoggerConfigurations))

	// create the loggerManager
	glog.Info("Create logger manager")
	loggerManager = log.NewLoggerManager(config.LoggerConfigurations)
	go loggerManager.Run()
	defer loggerManager.Stop()

	gui := gui.NewGui(app, loggerManager)

	// ESC exits
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			gui.Stop()
			loggerManager.Stop()
			app.Stop()
		}
		gui.HandleEventKey(event)
		return event
	})

	app.SetRoot(gui.Layout(), true)
	gui.Start()
	app.Run()
}
