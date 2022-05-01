package main

import (
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
	er.GetRules()
}
