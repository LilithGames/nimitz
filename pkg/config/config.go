package config

import (
	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"log"
)

var cfg = viper.New()

// Initialize configuration from yaml
func Initialize() (*viper.Viper, error) {

	// Viper config management initialization
	cfg.SetConfigType("yaml")
	cfg.SetConfigName("config")
	//cfg.AddConfigPath(".")
	//cfg.AddConfigPath("config")
	//cfg.AddConfigPath("pkg/config")
	//cfg.AddConfigPath("nimitz/pkg/config")
	cfg.AddConfigPath("/etc/webhook/config")
	cfg.AutomaticEnv()

	err := cfg.ReadInConfig()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("unable to marshal config to json: %s\n", err)
	}

	cfg.WatchConfig()
	return cfg, err
}

// GetCfg provide single instance of viper
func GetCfg() *viper.Viper {
	return cfg
}

// ImageMap loading self-define images match rules form yaml
func ImageMap() ([]interface{}, error) {
	newCfg := GetCfg()
	var imageRules = make(map[interface{}]interface{})
	val := newCfg.Get("ImageRules")
	rulesList := make([]interface{}, 0)
	for _, v := range val.([]interface{}) {
		switch v.(type) {
		case map[interface{}]interface{}:
			x := v.(map[interface{}]interface{})
			templist := make([]string, 2)
			templist[0] = x["pattern"].(string)
			templist[1] = x["replace"].(string)
			rulesList = append(rulesList, templist)
			imageRules[x["pattern"].(string)] = x["replace"].(string)
		default:
			log.Fatalf("%v\n", v)
		}
	}
	return rulesList, nil
}
