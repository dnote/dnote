package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
)

type Config struct {
	Channel string
}

type Note struct {
	Channel string
}

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.dnote", usr.HomeDir), nil
}

func touchDotDnote() error {
	content := []byte("channel: general\n")
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, content, 0644)
	return err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readConfig() (Config, error) {
	var ret Config

	configPath, err := getConfigPath()
	if err != nil {
		return ret, err
	}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ret, err
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func writeConfig(config Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// changeChannel replaces the channel name in the dnote config file
func changeChannel(channelName string) error {
	config, err := readConfig()
	if err != nil {
		return err
	}

	config.Channel = channelName

	err = writeConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if _, err := os.Stat("~/.dnote"); os.IsNotExist(err) {
		err = touchDotDnote()
		check(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Dnote - A command line tool to spontaneously record new learnings\n")
		fmt.Println("Commands:")
		fmt.Println("  use - choose the channel")
		fmt.Println("  new - write a new note")
		fmt.Println("  channels - show channels")
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "use":
		channel := os.Args[2]
		err := changeChannel(channel)
		check(err)
	}
}
