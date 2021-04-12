package config

// For the sake of simplicity, we will use one struct for both our API server and worker.
// To be able to read the config from both a local file as well as environment variables,
// we will use the library Viper, https://github.com/spf13/viper

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type CadenceConfig struct {
	Domain   string
	Service  string
	HostPort string
}

type AppConfig struct {
	Env            string
	WorkerTaskList string
	Cadence        CadenceConfig
	Logger         *zap.Logger
}

// Setup setup the config for the code run
func (h *AppConfig) Setup() {
	viper.SetConfigName("application")
	viper.AddConfigPath("app/resources") // These two lines will make sure viper pulls the config from app/resources/application.yml
	viper.AutomaticEnv() // This allows viper to read variables from the environment variables if they exists.
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}


	err := viper.Unmarshal(&h)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	hostPort := viper.GetString(h.Cadence.HostPort)
	fmt.Println(hostPort)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	h.Logger = logger

	logger.Debug("Finished loading Configuration!")
}