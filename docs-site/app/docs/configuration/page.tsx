import { DocPagination } from "@/components/doc-pagination"

export default function ConfigurationPage() {
  return (
    <div className="space-y-6">
      <h1>Configuration</h1>
      <p>
        GoFrame uses a simple YAML-based configuration system. All configuration settings are stored in the
        <code>config.yaml</code> file in the root of your project.
      </p>

      <h2>Basic Configuration</h2>
      <p>
        The <code>config.yaml</code> file contains all the configuration settings for your application, including
        database connection details, server settings, authentication settings, and more.
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

      <h2>Loading Configuration</h2>
      <p>
        GoFrame automatically loads the configuration from the <code>config.yaml</code> file when your application
        starts. You can access the configuration values through the <code>config.Config</code> struct.
      </p>

      <pre>
        <code>{`// main.go
// Load configuration
cfg, err := config.Load("config.yaml")
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}

// Access configuration values
fmt.Println(cfg.Server.Host)
fmt.Println(cfg.Database.Name)
fmt.Println(cfg.Auth.Secret)`}</code>
      </pre>

      <h2>Environment Variables</h2>
      <p>
        You can use environment variables to override configuration values. This is useful for deploying your
        application to different environments.
      </p>

      <pre>
        <code>{`// config/config.go
func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    // Override with environment variables
    if host := os.Getenv("DB_HOST"); host != "" {
        config.Database.Host = host
    }
    if port := os.Getenv("DB_PORT"); port != "" {
        portInt, _ := strconv.Atoi(port)
        config.Database.Port = portInt
    }
    if name := os.Getenv("DB_NAME"); name != "" {
        config.Database.Name = name
    }
    if user := os.Getenv("DB_USER"); user != "" {
        config.Database.User = user
    }
    if password := os.Getenv("DB_PASSWORD"); password != "" {
        config.Database.Password = password
    }

    return &config, nil
}`}</code>
      </pre>

      <h2>Configuration Structure</h2>
      <p>
        The configuration is structured as a nested set of structs. Here's the definition of the <code>Config</code>
        struct:
      </p>

      <pre>
        <code>{`// config/config.go
type Config struct {
    Server struct {
        Host string \`yaml:"host"\`
        Port int    \`yaml:"port"\`
    } \`yaml:"server"\`
    Database struct {
        Driver   string \`yaml:"driver"\`
        Host     string \`yaml:"host"\`
        Port     int    \`yaml:"port"\`
        Name     string \`yaml:"name"\`
        User     string \`yaml:"user"\`
        Password string \`yaml:"password"\`
    } \`yaml:"database"\`
    Auth struct {
        Secret   string        \`yaml:"secret"\`
        Duration time.Duration \`yaml:"duration"\`
    } \`yaml:"auth"\`
    RateLimit struct {
        Requests int           \`yaml:"requests"\`
        Period   time.Duration \`yaml:"period"\`
    } \`yaml:"rateLimit"\`
}`}</code>
      </pre>

      <h2>Custom Configuration</h2>
      <p>
        You can extend the configuration structure to include your own custom settings. Simply add new fields to the
        <code>Config</code> struct and update the <code>config.yaml</code> file accordingly.
      </p>

      <pre>
        <code>{`// config/config.go
type Config struct {
    // Existing fields...
    
    // Custom settings
    App struct {
        Name        string \`yaml:"name"\`
        Environment string \`yaml:"environment"\`
        Debug       bool   \`yaml:"debug"\`
    } \`yaml:"app"\`
}

// config.yaml
app:
  name: My GoFrame App
  environment: development
  debug: true`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Project Structure",
          href: "/docs/project-structure",
        }}
        next={{
          title: "Core Concepts",
          href: "/docs/core-concepts",
        }}
      />
    </div>
  )
}

