package env

import (
	"fmt"
	"github.com/spf13/viper"
)


func LoadEnvironment() {

	// Set the file name of the configurations file
	viper.SetConfigName("env")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		panic(err)
	}

	// Set undefined variables
	viper.SetDefault("app.debug", true)
	viper.SetDefault("app.bcryptCost", 10)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.maxOpenConnections", 20)
	viper.SetDefault("database.maxIdleConnections", 20)
	viper.SetDefault("database.maxLifetime", 300)
	viper.SetDefault("database.autoMigrate", false)

}
