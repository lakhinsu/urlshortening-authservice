package utils

import (
	"fmt"

	"github.com/lakhinsu/urlshortening-authservice/consts"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(consts.ENV_FILE)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("warning config file: %w", err))
	}
	// TODO: Implement this to work from os.env as well
	viper.AutomaticEnv()
}

func ReadEnvVar(name string) string {
	if !viper.IsSet(name) {
		msg := fmt.Sprintf(`Environment varaible %s is not set`, name)
		fmt.Println(msg)
		return ""
	}
	value := viper.GetString(name)
	return value
}
