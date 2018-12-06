package commands

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"net"
	"regexp"
)

var ConfigFilename = "config.json"

type Configuration struct {
	ServerAddr string
	ServerPort int
	Password string
}

func (x *RunCommand) Execute(args []string) error {
	var err error
	var config Configuration

	err, config = loadConfig()
	if err != nil {
		fmt.Print("Input server hostname: ")
		_, err = fmt.Scanf("%s\n", &config.ServerAddr)
		checkError(err)

		fmt.Print("Input server tunnel port: ")
		_, err = fmt.Scanf("%d\n", &config.ServerPort)
		checkError(err)

		fmt.Print("Input access password: ")
		_, err = fmt.Scanf("%s\n", &config.Password)
		checkError(err)

		// check if new configuration is valid
		err = checkConfig(config)
		checkError(err)

		// if configuration is valid - save data to file
		configJSON, _ := json.Marshal(config)
		err = ioutil.WriteFile(ConfigFilename, configJSON, 0644)
		checkError(err)
	}

	err = checkConfig(config)
	checkError(err)

	return nil
}

func loadConfig() (error, Configuration) {
	plainText, _ := ioutil.ReadFile(ConfigFilename)
	var config = Configuration{}
	err := json.Unmarshal(plainText, &config)

	return err, config
}

func checkConfig(config Configuration) error {
	_, err := net.LookupHost(config.ServerAddr)
	if err != nil {
		return err
	}

	if config.ServerPort < 1024 || config.ServerPort > 65535 {
		return errors.New("Invalid server port. Should be in the range from 1024 to 65535")
	}

	match, _ := regexp.MatchString("^[a-zA-Z0-9+/]+[=]*$", config.Password)
	if !match {
		return errors.New("Invalid password. Should be a valid string in base64 format")
	}

	return nil
}