package awsshield

import (
	"github.com/alecthomas/kingpin"
	"time"
)

type Config struct {
	Interval time.Duration
}

var defaultConfig = &Config{
	Interval: time.Minute,
}

func NewConfig() *Config {
	return &Config{}
}

// ParseFlags adds and parses flags from command line
func (cfg *Config) ParseFlags(args []string) error {
	app := kingpin.New("aws-shield-controller", "AWS Shield Controller turns on AWS Shield on AWS Load Balancers")
	app.Version("unkown")
	app.DefaultEnvars()
	app.Flag("interval", "The interval between two consecutive synchronizations in duration format (default: 1m)").Default(defaultConfig.Interval.String()).DurationVar(&cfg.Interval)

	_, err := app.Parse(args)
	if err != nil {
		return err
	}

	return nil
}
