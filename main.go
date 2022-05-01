package main

import (
	"fmt"
	"log"

	"http-testing/pkg/easyredir"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	APIKey    string `env:"EASYREDIR_API_KEY"`
	APISecret string `env:"EASYREDIR_API_SECRET"`
}

func main() {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatal("unable to read configuration")
	}

	er := easyredir.New(
		easyredir.WithAPIKey(cfg.APIKey),
		easyredir.WithAPISecret(cfg.APISecret),
	)

	rules, err := er.GetRules()
	if err != nil {
		log.Fatal("unable to get rules")
	}

	fmt.Print(rules)
}
