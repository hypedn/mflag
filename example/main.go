package main

import (
	"fmt"
	"log"

	"github.com/hypedn/mflag"
)

type AppSettings struct {
	Debug    bool
	AppPort  int
	Database struct {
		Host string
		Port int
		User string
	}
	Features []string
}

func GetSettings() AppSettings {
	return AppSettings{
		Debug:    mflag.GetBool(debug),
		AppPort:  mflag.GetInt(appPort),
		Features: mflag.GetStringSlice("features"),
		Database: struct {
			Host string
			Port int
			User string
		}{
			Host: mflag.GetString(database + "." + databaseHost),
			Port: mflag.GetInt(database + "." + databasePort),
			User: mflag.GetString(database + "." + databaseUser),
		},
	}
}

const (
	debug   = "debug"
	appPort = "app_port"

	database     = "database"
	databaseHost = "host"
	databasePort = "port"
	databaseUser = "user"

	features    = "features"
	darkMode    = "dark_mode"
	betaTesting = "beta_testing"
)

func defaults() {
	mflag.SetDefault(debug, true)
	mflag.SetDefault(appPort, 8080)

	mflag.SetDefault(database, map[string]interface{}{
		databaseHost: "localhost",
		databasePort: 5432,
		databaseUser: "default_user",
	})

	mflag.SetDefault(features, []string{darkMode, betaTesting})
}

func main() {
	defaults()
	if err := mflag.Init("configmap.yaml"); err != nil {
		log.Fatal(err)
	}
	dbHost := mflag.String("database.host", "localhost", "Database host address")
	appPort := mflag.Int("app_port", 3000, "Server port")
	debug := mflag.Bool("debug", true, "Enable debug mode")
	mflag.Parse()

	config := GetSettings()

	fmt.Println("Configuration loaded successfully!")
	fmt.Println("=================================")

	// The Get* functions now reflect the final merged configuration.
	fmt.Printf("Server Port: %d\n", config.AppPort)
	fmt.Printf("Debug Mode: %t\n", config.Debug)
	fmt.Printf("Database Host: %s\n", config.Database.Host)
	fmt.Printf("Database Port: %d\n", config.Database.Port)
	fmt.Printf("Database User: %s\n", config.Database.User)
	fmt.Printf("Enabled Features: %v\n", config.Features)

	// Using flag variables directly
	fmt.Printf("\nUsing flag variables directly:\n")
	fmt.Printf("Direct DB Host: %s\n", *dbHost)
	fmt.Printf("Direct App Port: %d\n", *appPort)
	fmt.Printf("Direct Debug: %t\n", *debug)

	fmt.Println("\n--- For Debugging ---")
	mflag.Debug()
}
