// This program performs administrative tasks for the garage sale service.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexlucaci/cookie-syncing-services/app/admin/commands"
	"github.com/alexlucaci/cookie-syncing-services/foundation/database"
	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		if errors.Cause(err) != commands.ErrHelp {
			log.Printf("error: %s", err)
		}
		os.Exit(1)
	}
}

func run() error {

	// =========================================================================
	// Configuration

	var cfg struct {
		conf.Version
		Args conf.Args
		DB   struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:cookie"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	const prefix = "Admin"
	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)

	// =========================================================================
	// Commands

	dbConfig := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	switch cfg.Args.Num(0) {
	case "migrate":
		if err := commands.Migrate(dbConfig); err != nil {
			return errors.Wrap(err, "migrating database")
		}

	default:
		fmt.Println("migrate: create the schema in the database")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}
