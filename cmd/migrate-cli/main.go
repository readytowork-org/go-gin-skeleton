package main

import (
	"boilerplate-api/bootstrap"
	"boilerplate-api/lib/config"
	"flag"
	"log"
	"os"

	"go.uber.org/fx"
)

// migrateApp provides just the dependencies needed for migration
var migrateModule = fx.Options(
	bootstrap.Module,
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("Usage: migrate-cli [up|down]")
	}

	command := args[0]

	// create app with dependencies
	app := fx.New(
		migrateModule,
		fx.Invoke(func(migrations *config.Migrations) {
			var err error
			switch command {
			case "up":
				err = migrations.MigrateUp()
			case "down":
				err = migrations.MigrateDown()
			default:
				log.Fatalf("Unknown command: %s", command)
			}
			if err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			os.Exit(0)
		}),
	)

	app.Run()
}
