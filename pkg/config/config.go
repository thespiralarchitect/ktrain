package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func BindDefault(serviceName string) error {
	filePath := fmt.Sprintf("../%s/config.yaml", serviceName)
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%v is not a file", filePath)
	}

	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")

	return viper.MergeInConfig()
}
