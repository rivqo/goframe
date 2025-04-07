package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/goframe/cmd/goframe/commands"
	"github.com/example/goframe/config"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Parse command
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "migrate":
		handleMigrate(cfg, args)
	case "make:migration":
		handleMakeMigration(args)
	case "make:model":
		handleMakeModel(args)
	case "make:controller":
		handleMakeController(args)
	case "make:resource":
		handleMakeResource(args)
	case "serve":
		handleServe(cfg)
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func handleMigrate(cfg *config.Config, args []string) {
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	rollback := migrateCmd.Bool("rollback", false, "Rollback the last migration batch")
	reset := migrateCmd.Bool("reset", false, "Rollback all migrations")
	refresh := migrateCmd.Bool("refresh", false, "Rollback all migrations and run them again")
	step := migrateCmd.Int("step", 0, "Number of migrations to rollback")

	migrateCmd.Parse(args)

	if *rollback {
		commands.MigrateRollback(cfg, *step)
	} else if *reset {
		commands.MigrateReset(cfg)
	} else if *refresh {
		commands.MigrateRefresh(cfg)
	} else {
		commands.Migrate(cfg)
	}
}

func handleMakeMigration(args []string) {
	if len(args) < 1 {
		fmt.Println("Migration name is required")
		os.Exit(1)
	}
	name := args[0]
	commands.MakeMigration(name)
}

func handleMakeModel(args []string) {
	if len(args) < 1 {
		fmt.Println("Model name is required")
		os.Exit(1)
	}
	name := args[0]
	migration := false
	if len(args) > 1 && args[1] == "-m" {
		migration = true
	}
	commands.MakeModel(name, migration)
}

func handleMakeController(args []string) {
	if len(args) < 1 {
		fmt.Println("Controller name is required")
		os.Exit(1)
	}
	name := args[0]
	resource := false
	if len(args) > 1 && args[1] == "-r" {
		resource = true
	}
	commands.MakeController(name, resource)
}

func handleMakeResource(args []string) {
	if len(args) < 1 {
		fmt.Println("Resource name is required")
		os.Exit(1)
	}
	name := args[0]
	commands.MakeResource(name)
}

func handleServe(cfg *config.Config) {
	commands.Serve(cfg)
}

func printHelp() {
	fmt.Println("GoFrame CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  goframe [command] [options]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  migrate                Run database migrations")
	fmt.Println("  migrate --rollback     Rollback the last migration batch")
	fmt.Println("  migrate --reset        Rollback all migrations")
	fmt.Println("  migrate --refresh      Rollback all migrations and run them again")
	fmt.Println("  migrate --step=n       Rollback n migrations")
	fmt.Println("  make:migration [name]  Create a new migration file")
	fmt.Println("  make:model [name]      Create a new model")
	fmt.Println("  make:controller [name] Create a new controller")
	fmt.Println("  make:resource [name]   Create a new resource")
	fmt.Println("  serve                  Start the HTTP server")
	fmt.Println("  help                   Display this help message")
}

