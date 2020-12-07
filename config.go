package main

import (
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("user", "u", "", "SSH Username")
	pflag.Bool("private", false, "Use Private IP address when connecting")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal("Couldn't bind flags")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ec2_fuzzy")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
