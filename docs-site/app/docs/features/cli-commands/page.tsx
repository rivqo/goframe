import { DocPagination } from "@/components/doc-pagination"

export default function CliCommandsPage() {
  return (
    <div className="space-y-6">
      <h1>CLI Commands</h1>
      <p>
        GoFrame provides a powerful command-line interface (CLI) that helps you manage your application. The CLI
        includes commands for common tasks like running migrations, creating models, controllers, and resources, and
        starting the development server.
      </p>

      <h2>Available Commands</h2>

      <h3>Server Commands</h3>
      <pre>
        <code>{`# Start the development server
goframe serve`}</code>
      </pre>

      <h3>Migration Commands</h3>
      <pre>
        <code>{`# Create a new migration
goframe make:migration create_users_table

# Run all pending migrations
goframe migrate

# Rollback the last batch of migrations
goframe migrate --rollback

# Rollback all migrations
goframe migrate --reset

# Rollback and re-run all migrations
goframe migrate --refresh

# Rollback a specific number of migrations
goframe migrate --step=2`}</code>
      </pre>

      <h3>Model Commands</h3>
      <pre>
        <code>{`# Create a new model
goframe make:model User

# Create a model with a migration
goframe make:model User -m`}</code>
      </pre>

      <h3>Controller Commands</h3>
      <pre>
        <code>{`# Create a new controller
goframe make:controller UserController

# Create a resource controller
goframe make:controller PostController -r`}</code>
      </pre>

      <h3>Resource Commands</h3>
      <pre>
        <code>{`# Create a new resource
goframe make:resource Post`}</code>
      </pre>

      <h2>Creating Custom Commands</h2>
      <p>
        You can create your own custom commands by adding them to the <code>cmd/goframe/commands</code> directory:
      </p>
      <pre>
        <code>{`// cmd/goframe/commands/custom.go
package commands

import (
    "fmt"
    "os"
)

// CustomCommand handles a custom command
func CustomCommand(args []string) {
    fmt.Println("Running custom command...")
    
    // Your command logic here
    
    fmt.Println("Custom command completed successfully!")
}
`}</code>
      </pre>

      <p>
        Then register your command in the <code>cmd/goframe/main.go</code> file:
      </p>
      <pre>
        <code>{`// cmd/goframe/main.go
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
    case "custom": // Add your custom command here
        commands.CustomCommand(args)
    case "help":
        printHelp()
    default:
        fmt.Printf("Unknown command: %s\n", command)
        printHelp()
        os.Exit(1)
    }
}`}</code>
      </pre>

      <h2>Command Structure</h2>
      <p>GoFrame commands follow a simple structure:</p>
      <pre>
        <code>{`goframe [command] [options]`}</code>
      </pre>

      <p>For example:</p>
      <pre>
        <code>{`goframe make:migration create_users_table
goframe migrate --rollback
goframe make:model User -m`}</code>
      </pre>

      <h2>Getting Help</h2>
      <p>You can get help on available commands by running:</p>
      <pre>
        <code>{`goframe help`}</code>
      </pre>

      <p>This will display a list of all available commands and their options.</p>

      <DocPagination
        prev={{
          title: "Rate Limiting",
          href: "/docs/features/rate-limiting",
        }}
        next={{
          title: "Advanced",
          href: "/docs/advanced",
        }}
      />
    </div>
  )
}

