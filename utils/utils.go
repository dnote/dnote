package utils

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"strconv"
	"time"
)

const configFilename = ".dnoterc"
const DnoteUpdateFilename = ".dnote-update"
const dnoteFilename = ".dnote"

func GetConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, configFilename), nil
}

func GetDnotePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, dnoteFilename), nil
}

func GenerateConfigFile() error {
	content := []byte("book: general\n")
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, content, 0644)
	return err
}

func TouchDnoteFile() error {
	dnotePath, err := GetDnotePath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dnotePath, []byte{}, 0644)
	return err
}

func TouchDnoteUpgradeFile() error {
	fmt.Println("toching dnoteupdate")
	dnoteUpdatePath, err := GetDnoteUpdatePath()
	if err != nil {
		return err
	}

	epoch := strconv.FormatInt(time.Now().Unix(), 10)
	content := []byte(fmt.Sprintf("LAST_UPGRADE_EPOCH: %s\n", epoch))

	err = ioutil.WriteFile(dnoteUpdatePath, content, 0644)
	return err
}

func GetDnoteUpdatePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, DnoteUpdateFilename), nil
}
