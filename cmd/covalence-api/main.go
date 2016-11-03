package main

import (
	"fmt"
	"os"

	"github.com/starkandwayne/covalence/api"
	"github.com/starkandwayne/goutils/log"
	"github.com/voxelbrain/goptions"

	// sql drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type CovalanceServerOpts struct {
	Help       bool   `goptions:"-h, --help, description='Show the help screen'"`
	ConfigFile string `goptions:"-c, --config, description='Path to the api configuration file'"`
	Log        string `goptions:"-l, --log-level, description='Set logging level to debug, info, notice, warn, error, crit, alert, or emerg'"`
	Version    bool   `goptions:"-v, --version, description='Display the api version'"`
}

var Version = ""

func main() {
	api.Version = Version
	var opts CovalanceServerOpts
	opts.Log = "Info"
	if err := goptions.Parse(&opts); err != nil {
		fmt.Printf("%s\n", err)
		goptions.PrintHelp()
		return
	}

	if opts.Help {
		goptions.PrintHelp()
		os.Exit(0)
	}
	if opts.Version {
		if Version == "" {
			fmt.Printf("covalence-api (development)\n")
		} else {
			fmt.Printf("covalence-api v%s\n", Version)
		}
		os.Exit(0)
	}

	if opts.ConfigFile == "" {
		fmt.Fprintf(os.Stderr, "No config specified. Please try again using the -c/--config argument\n")
		os.Exit(1)
	}

	log.SetupLogging(log.LogConfig{Type: "console", Level: opts.Log})
	log.Infof("starting covalence api")

	a := api.NewApi()
	if err := a.ReadConfig(opts.ConfigFile); err != nil {
		log.Errorf("Failed to load config: %s", err)
		return
	}

	if err := a.Database.Connect(); err != nil {
		log.Errorf("failed to connect to %s database at %s: %s\n",
			a.Database.Driver, a.Database.DSN, err)
	}

	if err := a.Database.Setup(); err != nil {
		log.Errorf("failed to set up schema in %s database at %s: %s\n",
			a.Database.Driver, a.Database.DSN, err)
	}

	if err := a.Run(); err != nil {
		log.Errorf("covalence api failed: %s", err)
	}
	log.Infof("stopping covalence api")
}
