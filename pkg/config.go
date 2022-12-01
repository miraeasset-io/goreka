package goreka

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	RegistryType    string
	EurekaUrl       string
	EurekaUser      string
	EurekaPassword  string
	HotsAddressPort string
	ServicePort     int
	ServiceName     string
	InstanceId      string
	AppProfile      string
	ServiceHost     string
}

func LoadConfig() (config Config) {
	viper.SetEnvPrefix("EUREKA")
	viper.AutomaticEnv()

	port, err := strconv.Atoi(fmt.Sprintf("%v", viper.Get("SERVICE_PORT")))
	if err != nil {
		panic(err)
	}

	return Config{
		RegistryType:    fmt.Sprintf("%v", viper.Get("REGISTRY_TYPE")),
		EurekaUrl:       fmt.Sprintf("%v", viper.Get("URL")),
		EurekaUser:      fmt.Sprintf("%v", viper.Get("USERNAME")),
		EurekaPassword:  fmt.Sprintf("%v", viper.Get("PASSWORD")),
		HotsAddressPort: fmt.Sprintf("%v", viper.Get("HOTS_ADDRESS_PORT")),
		ServiceHost:     fmt.Sprintf("%v", viper.Get("SERVICE_HOST")),
		ServicePort:     port,
		ServiceName:     fmt.Sprintf("%v", viper.Get("SERVICE_NAME")),
		InstanceId:      fmt.Sprintf("%v", viper.Get("INSTANCE_ID")),
		AppProfile:      fmt.Sprintf("%v", viper.Get("APP_PROFILE")),
	}
}
