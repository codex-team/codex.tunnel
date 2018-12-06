package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
)

var ConfigFilename = "config.json"

type Configuration struct {
	ServerAddr string
	ServerPort int
	Password string
	PrivateKeyPath string
}

type RegistrationMessage struct {
	Key string `json:"key"`
	Password string `json:"password"`
}

func (x *RunCommand) Execute(args []string) error {
	var err error
	var config Configuration

	err, config = loadConfig()
	if err != nil {
		fmt.Print("Input server address (ex: https://tun.ifmo.su or http://localhost:8081): ")
		_, err = fmt.Scanf("%s\n", &config.ServerAddr)
		checkError(err)

		fmt.Print("Input server tunnel port (ex: 17022): ")
		_, err = fmt.Scanf("%d\n", &config.ServerPort)
		checkError(err)

		fmt.Print("Input access password: ")
		_, err = fmt.Scanf("%s\n", &config.Password)
		checkError(err)

		// check if new configuration is valid
		err = checkConfig(config)
		checkError(err)

		keyName := string(RandASCIIBytes(8))
		privKeyPath, SshKey, err := generateKeySet(keyName)
		checkAndClear(keyName, err)

		config.PrivateKeyPath = privKeyPath

		// try to publish register key
		err = registerKey(RegistrationMessage{SshKey, config.Password}, config)
		checkAndClear(keyName, err)

		// if configuration is valid and key is successfully published to the server - save data to a file
		configJSON, _ := json.Marshal(config)
		err = ioutil.WriteFile(ConfigFilename, configJSON, 0644)
		checkAndClear(keyName, err)
	}

	err = checkConfig(config)
	checkError(err)

	host := x.Host
	if host == "" {
		host = string(RandASCIIBytes(6))
	}

	return make_tunnel(host, x.LocalHost, x.Port, config.ServerAddr, config.ServerPort, config.PrivateKeyPath)
}

func loadConfig() (error, Configuration) {
	var config = Configuration{}
	plainText, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		return err, config
	}
	err = json.Unmarshal(plainText, &config)

	return err, config
}

func checkConfig(config Configuration) error {
	u, err := url.Parse(config.ServerAddr)
	if err != nil {
		return err
	}

	_, err = net.LookupHost(u.Hostname())
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

func registerKey(message RegistrationMessage, config Configuration) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(message)

	u, err := url.Parse(config.ServerAddr)
	u.Path = path.Join(u.Path, "/register")

	res, err := http.Post(u.String(), "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}

	responseText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(string(responseText))
	}

	return nil
}

func generateKeySet(name string) (string, string, error) {
	privKeyPath := fmt.Sprintf("./keys/%s_priv.pem", name)
	pubKeyPath := fmt.Sprintf("./keys/%s_pub.pem", name)
	sshKeyPath := fmt.Sprintf("./keys/%s.key", name)

	err := generateKeys(privKeyPath, pubKeyPath, sshKeyPath)
	if err != nil {
		return "", "", err
	}

	sshKey, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return "", "", err
	}

	return privKeyPath, string(sshKey)[8:], nil
}

func checkAndClear(name string, err error)  {
	if err != nil {
		clear(name)
	}
}

func clear(name string) {
	privKeyPath := fmt.Sprintf("./keys/%s_priv.pem", name)
	pubKeyPath := fmt.Sprintf("./keys/%s_pub.pem", name)
	sshKeyPath := fmt.Sprintf("./keys/%s.key", name)
	os.Remove(privKeyPath)
	os.Remove(pubKeyPath)
	os.Remove(sshKeyPath)
}