package main

import (
	"kustocheck/files"
	"log"
)

var configFile = "../config/config.yaml"

func main() {
	_, err := files.GetPats(configFile)

	if err != nil {
		log.Fatal(err)
	}

}
