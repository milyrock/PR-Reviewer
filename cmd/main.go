package main

import (
	"fmt"

	"github.com/milyrock/PR-Reviewer/internal/config"
)

func main() {
	config, err := config.ReadConfig("./config/config.yaml")
	if err != nil {
		fmt.Println("Couldn't read the config", err)
	}
	fmt.Println(config)
}
