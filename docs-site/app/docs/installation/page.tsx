import { DocPagination } from "@/components/doc-pagination"

export default function InstallationPage() {
  return (
    <div className="space-y-6">
      <h1>Installation</h1>
      <p>This guide will help you install GoFrame and set up your first project.</p>
      <h2>Prerequisites</h2>
      <p>Before you begin, make sure you have the following installed:</p>
      <ul>
        <li>Go 1.18 or higher</li>
        <li>PostgreSQL (or another supported database)</li>
      </ul>
      <h2>Installing GoFrame</h2>
      <p>You can install GoFrame by cloning the repository and building the CLI tool:</p>
      <pre>
        <code>{`# Clone the repository
git clone https://github.com/yourusername/goframe.git

# Navigate to the project directory
cd goframe

# Install dependencies
go mod tidy

# Build the CLI tool
go build -o goframe ./cmd/goframe

# Make the CLI tool available globally (optional)
# On Linux/macOS:
sudo mv goframe /usr/local/bin/
# On Windows:
# Move goframe.exe to a directory in your PATH`}</code>
      </pre>
      <h2>Creating a New Project</h2>
      <p>Once you have the GoFrame CLI installed, you can create a new project:</p>
      <pre>
        <code>{`# Create a new project directory
mkdir myproject
cd myproject

# Initialize a Go module
go mod init github.com/yourusername/myproject

# Initialize the project structure
goframe init`}</code>
      </pre>
      <h2>Project Structure</h2>
      <p>After initializing your project, you'll have the following structure:</p>
      <pre>
        <code>{`myproject/
├── cmd/                  # Command-line applications
├── config/               # Configuration files
├── controllers/          # HTTP controllers
├── db/                   # Database connection and utilities
├── middleware/           # HTTP middleware
├── migrations/           # Database migrations
├── models/               # Data models
├── resources/            # API resources
├── routes/               # Route definitions
├── config.yaml           # Application configuration
├── go.mod                # Go module file
├── go.sum                # Go module checksum
└── main.go               # Application entry point`}</code>
      </pre>
      <h2>Configuration</h2>
      <p>
        Configure your application in <code>config.yaml</code>:
      </p>
      <pre>
        <code>{`server:
  host: localhost
  port: 8080

database:
  driver: postgres
  host: localhost
  port: 5432
  name: goframe
  user: postgres
  password: postgres

auth:
  secret: your-secret-key-here
  duration: 24h

rateLimit:
  requests: 100
  period: 1m`}</code>
      </pre>
      <h2>Starting the Server</h2>
      <p>To start the development server, run:</p>
      <pre>
        <code>{`goframe serve`}</code>
      </pre>
      <p>
        Your application will be available at <code>http://localhost:8080</code>.
      </p>
      <h2>Next Steps</h2>
      <p>
        Now that you have GoFrame installed, check out the <a href="/docs/project-structure">Project Structure</a> guide
        to learn more about how GoFrame projects are organized.
      </p>
      <DocPagination
        prev={{
          title: "Introduction",
          href: "/docs",
        }}
        next={{
          title: "Project Structure",
          href: "/docs/project-structure",
        }}
      />
    </div>
  )
}

