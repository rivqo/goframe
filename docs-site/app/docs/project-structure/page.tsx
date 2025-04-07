import { DocPagination } from "@/components/doc-pagination"

export default function ProjectStructurePage() {
  return (
    <div className="space-y-6">
      <h1>Project Structure</h1>
      <p>
        GoFrame follows a clean and organized project structure that makes it easy to find and manage your code. This
        page describes the default project structure and explains the purpose of each directory.
      </p>

      <h2>Directory Structure</h2>
      <pre>
        <code>{`myproject/
├── cmd/                  # Command-line applications
│   └── goframe/          # GoFrame CLI tool
│       ├── commands/     # CLI commands
│       └── main.go       # CLI entry point
├── config/               # Configuration files
│   └── config.go         # Configuration loader
├── controllers/          # HTTP controllers
│   ├── auth_controller.go
│   ├── user_controller.go
│   └── web_controller.go
├── db/                   # Database connection and utilities
│   ├── db.go             # Database connection
│   ├── migrator.go       # Migration system
│   └── orm.go            # ORM functionality
├── middleware/           # HTTP middleware
│   └── middleware.go     # Common middleware
├── migrations/           # Database migrations
│   ├── 20230615120000_create_users_table.go
│   └── ...
├── models/               # Data models
│   ├── user.go
│   └── ...
├── public/               # Public assets
│   ├── assets/           # Static assets
│   │   ├── css/          # CSS files
│   │   ├── js/           # JavaScript files
│   │   └── images/       # Image files
│   ├── favicon.ico       # Favicon
│   └── robots.txt        # Robots file
├── resources/            # API resources
│   ├── resource.go
│   └── user_resource.go
├── routes/               # Route definitions
│   ├── api.go            # API routes
│   ├── router.go         # Router initialization
│   └── web.go            # Web routes
├── views/                # View templates
│   ├── layouts/          # Layout templates
│   │   └── app.html      # Main layout
│   ├── pages/            # Page templates
│   │   ├── home.html
│   │   └── ...
│   ├── partials/         # Partial templates
│   │   ├── header.html
│   │   └── footer.html
│   └── errors/           # Error templates
│       └── 404.html
├── config.yaml           # Application configuration
├── go.mod                # Go module file
├── go.sum                # Go module checksum
└── main.go               # Application entry point`}</code>
      </pre>

      <h2>Key Directories</h2>

      <h3>cmd/</h3>
      <p>
        The <code>cmd/</code> directory contains command-line applications, including the GoFrame CLI tool. This is
        where you'll find the code for commands like <code>migrate</code>, <code>make:model</code>, etc.
      </p>

      <h3>config/</h3>
      <p>
        The <code>config/</code> directory contains configuration-related code, including the configuration loader that
        reads settings from <code>config.yaml</code>.
      </p>

      <h3>controllers/</h3>
      <p>
        The <code>controllers/</code> directory contains HTTP controllers that handle incoming requests and return
        responses. Controllers are organized by resource or functionality.
      </p>

      <h3>db/</h3>
      <p>
        The <code>db/</code> directory contains database-related code, including the connection manager, query builder,
        ORM functionality, and migration system.
      </p>

      <h3>middleware/</h3>
      <p>
        The <code>middleware/</code> directory contains HTTP middleware that can be applied to routes or groups of
        routes. Middleware can perform tasks like authentication, logging, rate limiting, etc.
      </p>

      <h3>migrations/</h3>
      <p>
        The <code>migrations/</code> directory contains database migrations that define your database schema. Migrations
        are versioned and can be run or rolled back as needed.
      </p>

      <h3>models/</h3>
      <p>
        The <code>models/</code> directory contains data models that represent tables in your database. Models typically
        include repository methods for interacting with the database.
      </p>

      <h3>public/</h3>
      <p>
        The <code>public/</code> directory contains publicly accessible files like CSS, JavaScript, images, and other
        assets. These files are served directly by the web server.
      </p>

      <h3>resources/</h3>
      <p>
        The <code>resources/</code> directory contains API resources that transform models into API responses. Resources
        help you build a consistent API layer.
      </p>

      <h3>routes/</h3>
      <p>
        The <code>routes/</code> directory contains route definitions that map URLs to controller methods. Routes are
        organized by type (web, API) and can be grouped and have middleware applied.
      </p>

      <h3>views/</h3>
      <p>
        The <code>views/</code> directory contains HTML templates for rendering web pages. Templates are organized by
        type (layouts, pages, partials) and can include each other.
      </p>

      <h2>Configuration Files</h2>

      <h3>config.yaml</h3>
      <p>
        The <code>config.yaml</code> file contains application configuration settings like database connection details,
        server settings, authentication settings, etc.
      </p>

      <h3>go.mod and go.sum</h3>
      <p>
        The <code>go.mod</code> and <code>go.sum</code> files are standard Go module files that define your
        application's dependencies and their versions.
      </p>

      <h3>main.go</h3>
      <p>
        The <code>main.go</code> file is the entry point for your application. It loads configuration, sets up the
        database connection, initializes the router, and starts the HTTP server.
      </p>

      <DocPagination
        prev={{
          title: "Installation",
          href: "/docs/installation",
        }}
        next={{
          title: "Configuration",
          href: "/docs/configuration",
        }}
      />
    </div>
  )
}

