package main

import (
  "fmt"
  "log"

	"http-testing/pkg/easyredir"
)

const (
	apiKey    string = "***REMOVED***"
	apiSecret string = "***REMOVED***"
)

func main() {
	er := easyredir.New(&easyredir.Options{
		APIKey:    apiKey,
		APISecret: apiSecret,
	})

	rules, err := er.GetRules()
  if err != nil {
    log.Fatal("unable to get rules")
  }

  fmt.Println(rules)
}
