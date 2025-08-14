# mflag

**Viper made simple with a flag-compatible API**

[![Go Reference](https://pkg.go.dev/badge/github.com/hypedn/mflag.svg)](https://pkg.go.dev/github.com/hypedn/mflag)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypedn/mflag)](https://goreportcard.com/report/github.com/hypedn/mflag)

mflag is a minimal configuration library for Go applications that combines the simplicity of the standard `flag` package with YAML configuration file support. If you love Viper's flexibility but want something lighter and more familiar, mflag is for you.

## ✨ Features

- **Drop-in replacement for `flag`** - Same API you already know
- **YAML configuration support** - Load defaults from config files
- **Clear precedence order** - Command-line flags > Config file > Code defaults
- **Minimal dependencies** - Only requires YAML parsing
- **Zero configuration** - Works out of the box with just flags

## 🚀 Quick Start

```bash
go get github.com/hypedn/mflag
```

```go
package main

import (
    "fmt"

    "github.com/hypedn/mflag"
)

func defaults() {
    mflag.SetDefault("port", 3000)
    mflag.SetDefault("host", "localhost")
    mflag.SetDefault("debug", false)
}

func main() {
    defaults() // Set default values
    mflag.Init("configmap.yaml") // Load configmap. Overrides defaults
    debug := mflag.Bool("debug", true, "Enable debug mode") // Set debug flag. Overrides configmap
    mflag.Parse() // Required to run Parse before reading configs.
    
    fmt.Printf("Server: %s:%d (debug: %t)\n", *host, *port, *debug)
}
```

**configmap.yaml:**
```yaml
port: 3000
host: "0.0.0.0"
debug: false
```

**Reading from yaml is optional and won't return an error if the file doesn't exist.**

You can also use flag with defaults, see [example](./example/main.go).

Values are resolved in this order (highest to lowest priority):

1. **Command-line flags** - Explicit user input
2. **YAML configuration file** - Persistent settings  
3. **Default values in code** - Fallback values

### Getting Values After Parsing

After calling `mflag.Parse()`, you can also retrieve values by key:

```go
mflag.Parse()

// Get typed values
port := mflag.GetInt("port")
host := mflag.GetString("host") 
debug := mflag.GetBool("debug")

// Check if a value was explicitly set
if mflag.IsSet("port") {
    fmt.Println("Port was explicitly configured")
}

// Get all configuration keys
for _, key := range mflag.AllKeys() {
    fmt.Printf("%s = %v\n", key, mflag.GetString(key))
}
```

### Setting Defaults Programmatically

A best practise is to set defaults before loading the config file.

```go
// Set defaults first
mflag.SetDefault("timeout", 30)
mflag.SetDefault("retries", 3)

// Then load config file
mflag.Init("config.yaml")

// Define flags
timeout := mflag.Int("timeout", 10, "Request timeout") // Will use 30 from SetDefault
// ... rest of your flags

mflag.Parse()
```

## 📋 Examples

### Basic Web Server

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/hypedn/mflag"
)

func main() {
    mflag.Init("server.yaml")
    
    port := mflag.Int("port", 8080, "Server port")
    host := mflag.String("host", "localhost", "Server host") 
    debug := mflag.Bool("debug", false, "Enable debug logging")
    
    mflag.Parse()
    
    if *debug {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    }
    
    addr := fmt.Sprintf("%s:%d", *host, *port)
    log.Printf("Starting server on %s", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}
```

### Database Configuration

**config.yaml:**
```yaml
database:
  host: "localhost"
  port: 5432
  name: "myapp"
  ssl_mode: "require"
max_connections: 100
timeout: 30
```

```go
func main() {
    mflag.Init("config.yaml")
    
    // Nested configuration with dot notation
    dbHost := mflag.String("database.host", "localhost", "Database host")
    dbPort := mflag.Int("database.port", 5432, "Database port")
    dbName := mflag.String("database.name", "app", "Database name")
    sslMode := mflag.String("database.ssl_mode", "disable", "SSL mode")
    
    // Top-level configuration
    maxConn := mflag.Int("max_connections", 10, "Max database connections")
    timeout := mflag.Int("timeout", 5, "Connection timeout (seconds)")
    
    mflag.Parse()
    
    dsn := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s", 
                      *dbHost, *dbPort, *dbName, *sslMode)
    fmt.Printf("DSN: %s\n", dsn)
    fmt.Printf("Max connections: %d, Timeout: %ds\n", *maxConn, *timeout)
}
```

### Debugging Configuration

```go
func main() {
    mflag.Init("config.yaml")
    
    // ... define your flags
    
    mflag.Parse()
    
    // Debug all configuration values
    mflag.Debug() // Prints all key-value pairs
}
```

## 🔄 Migration from flag

Migrating from the standard `flag` package is straightforward:

```go
// Before (using flag)
import "flag"

port := flag.Int("port", 8080, "Server port")
flag.Parse()

// After (using mflag)
import "github.com/hypedn/mflag"

mflag.Init("config.yaml") // Optional: add config file support
port := mflag.Int("port", 8080, "Server port")
mflag.Parse()
```

## 🔄 Migration from Viper

Coming from Viper? mflag offers a simpler approach:

```go
// Viper approach
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.ReadInConfig()
viper.SetDefault("port", 8080)
port := viper.GetInt("port")

// mflag approach  
mflag.Init("config.yaml")
port := mflag.Int("port", 8080, "Server port")
mflag.Parse()
```

## 🤔 Why mflag?

**Choose mflag when you want:**
- Familiar `flag` package API
- Simple YAML configuration support
- Minimal dependencies and overhead
- Clear, predictable behavior
- Easy migration from existing `flag` usage

**Stick with Viper when you need:**
- Multiple configuration formats (JSON, TOML, etc.)
- Environment variable support
- Remote configuration (etcd, Consul)
- Configuration watching/reloading
- Complex configuration transformations

## 📊 Comparison

| Feature | mflag | flag | Viper |
|---------|--------|------|-------|
| YAML config files | ✅ | ❌ | ✅ |
| Familiar API | ✅ | ✅ | ❌ |
| Minimal dependencies | ✅ | ✅ | ❌ |
| Multiple formats | ❌ | ❌ | ✅ |
| Environment variables | ❌ | ❌ | ✅ |
| Remote config | ❌ | ❌ | ✅ |

## 🛠️ Advanced Usage

### Custom YAML Structure

mflag supports nested YAML structures using dot notation:

```yaml
server:
  http:
    port: 8080
    host: "0.0.0.0"
  tls:
    enabled: true
    cert_file: "/path/to/cert.pem"
database:
  connection:
    host: "db.example.com"
    port: 5432
```

```go
httpPort := mflag.Int("server.http.port", 3000, "HTTP port")
tlsEnabled := mflag.Bool("server.tls.enabled", false, "Enable TLS")
dbHost := mflag.String("database.connection.host", "localhost", "Database host")
```

### Error Handling

```go
// Check if config file loading failed
if err := mflag.Init("config.yaml"); err != nil {
    log.Printf("Warning: Could not load config file: %v", err)
    // Continue with just flags and defaults
}

// ... define flags

mflag.Parse()
```

## 📝 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
