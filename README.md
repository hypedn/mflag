# mflag

[![Go Reference](https://pkg.go.dev/badge/github.com/hypedn/mflag.svg)](https://pkg.go.dev/github.com/hypedn/mflag)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypedn/mflag)](https://goreportcard.com/report/github.com/hypedn/mflag)

Go library that integrates config loading from YAML files and Go flags.

mflag was implemented based on a typical workflow for apps running in Kubernetes:
- In production, read configs from a configmap YAML file
- In staging, read default configs
- In local, read default configs with the option of overriding them via CLI flags, for ease of debugging.

## ✨ Features

- **Simple, Declarative API** - Define defaults, load a file, and parse. That's it.
- **Automatic Command-line Flags** - Every configuration key is automatically available as a command-line flag for overrides.
- **YAML configuration support** - Load defaults from config files
- **Clear precedence order** - Command-line flags > Config file > Code defaults
- **Minimal dependencies** - Only requires YAML parsing

## 🚀 Quick Start

```bash
go get github.com/hypedn/mflag
```

**main.go:**
```go
package main

import (
    "fmt"
    "log"

    "github.com/hypedn/mflag"
)

func defaults() {
    mflag.SetDefault("debug", true)
    mflag.SetDefault("port", 3000)
    mflag.SetDefault("host", "localhost")
    mflag.SetDefault("database.user", "default_user")
}

func main() {
    defaults()
    if err := mflag.Init("configmap.yaml"); err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    mflag.Parse()

    // Print all flags
    if mflag.GetBool("debug") {
        mflag.Debug()
    }
}
```

**configmap.yaml:**
```yaml
port: 3000
host: "0.0.0.0"
debug: false
database:
  host: "localhost"
  port: 5432
  name: "myapp"
  ssl_mode: "require"
```

You can override any config via CLI:
```bash
go run main.go --help
```

```bash
go run main.go --port=8080
```

## 📚 Good to know

**Reading from yaml is optional and won't return an error if the file doesn't exist**. Hence it is a good practise to always provide safe defaults.

Values are resolved in this order (highest to lowest priority):

1. **Command-line flags** - Explicit user input (highest priority)
2. **YAML configuration file** - Persistent settings
3. **Default values in code** - Fallback values

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

// Print all flags
mflag.Debug()
```

Check [example](./example/main.go) for a practical example of parsing configs into a struct. You May also want to split `AppConfig` into multiple configs like `DBConfig`, `CacheConfig`, etc.

## 🔧 Trade-offs

This library is for you if you want a lightweight and ergonomic API. This library is NOT for you if you require strong config validation at every step, and relying on sane defaults is not an option for your use case. See the examples below:

### Init() error handling

As mentioned above, we deliberately chose to not return an error when the config YAML file is not found. This allows an application to run with defaults locally without needing a config file. In any other case, such as a file with wrong permissions or invalid YAML syntax, `Init()` will return a descriptive error.

### Get* error handling

This library prioritizes ease of use over forcing error checks on every value retrieval. In the example above:

```go
port := mflag.GetInt("port")
host := mflag.GetString("host")
debug := mflag.GetBool("debug")
```

When the flag was not found, Get* functions return their zero value. If we would have opted for a stronger validation, we would need to check errors in every single line of loading each config, which would be verbose and tedious.

### Parse() error handling

`Parse()` does not return an error, to be consistent with the `flag` package from the standard library. It is designed for simplicity in common use cases. If an error occurs during parsing (e.g., an invalid value in a config file), `Parse()` will print the error message and exit the application, mirroring the default behavior of `flag.Parse()`.

Loading configurations is often the first thing an application does, so this is generally an acceptable approach. However if you open connections or spin up go routines before loading configs, exiting won't gracefully shutdown your application. For those cases you can use `ParseWithError()`. This function performs the same logic as Parse() but returns an error on failure instead of exiting.

Note: calling any Get* function before Parse() or ParseWithError() **will cause a panic**. This is a deliberate design choice to prevent silent failures from incorrect library usage, distinguishing a programmer error (violating the library's lifecycle) from a runtime error (bad input data).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
