package conf

import (
	"strings"
	"github.com/spf13/viper"
)

func TestInit(cfg string) error {
	c := Config {
		Name: cfg,
	}

	if err := c.testInitConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) testInitConfig() error  {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("../conf")
		viper.SetConfigName("conf")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO_DISTRIBUTE_STORAGE")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
