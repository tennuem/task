package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const ServiceName = "task"

var options = []option{
	{"config", "string", "", "config file"},

	{"server.http.host", "string", "0.0.0.0", "server grpc host"},
	{"server.http.port", "int", 8080, "server http port"},

	{"logger.level", "string", "emerg", "Level of logging. A string that correspond to the following levels: emerg, alert, crit, err, warning, notice, info, debug"},
	{"logger.time_format", "string", "2006-01-02T15:04:05.999999999", "Date format in logs"},
}

type Config struct {
	Server Server
	Logger Logger
}

type Server struct {
	HTTP HTTP
}

type HTTP struct {
	Host string
	Port int
}

type Logger struct {
	Level      string
	TimeFormat string
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}

// NewConfig returns and prints struct with config parameters
func NewConfig() *Config {
	cfg := &Config{}
	if err := cfg.read(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}
	if err := print(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "print config: %s", err)
		os.Exit(1)
	}
	return cfg
}

// read gets parameters from environment variables, flags or file.
func (c *Config) read() error {
	viper.SetEnvPrefix(ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, o := range options {
		switch o.typing {
		case "string":
			pflag.String(o.name, o.value.(string), o.description)
		case "int":
			pflag.Int(o.name, o.value.(int), o.description)
		default:
			viper.SetDefault(o.name, o.value)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigName(fileName)
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

func print(cfg *Config) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	log.Println(string(b))
	return nil
}
