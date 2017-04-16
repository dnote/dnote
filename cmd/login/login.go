package login

import (
	"fmt"

	"github.com/dnote-io/cli/utils"
)

func Run() error {
	fmt.Print("Please enter your APIKey: ")

	var apiKey string
	fmt.Scanln(&apiKey)

	config, err := utils.ReadConfig()
	if err != nil {
		return err
	}

	config.APIKey = apiKey
	err = utils.WriteConfig(config)
	if err != nil {
		return err
	}

	return nil
}
