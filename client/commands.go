package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var opts struct {
	Run RunCommand `command:"run" description:"Run tunel"`
	Generate GenerateCommand `command:"generate" description:"Generate rsa keys"`
	Tunnel TunnelCommand `command:"tunnel" description:"Start tunnel"`
}

var Configuration struct {
	ServerAddr string
	ServerPort int
	Password string
}

type TunnelCommand struct {
	Host string `short:"H" long:"hostname" description:"Desired hostname" required:"true"`
	LocalHost string `short:"L" long:"local_host" description:"Desired local host" required:"false" default:"localhost"`
	LocalPort int `short:"P" long:"local_port" description:"Desired local port" required:"true"`
	Key string `short:"k" long:"key" description:"Path to privkey file" required:"true"`
	ServerIp string `short:"h" long:"server_host" description:"Server hostname" required:"false" default:"tun.ifmo.su"`
	ServerPort int `short:"p" long:"server_port" description:"Server port" required:"false" default:"22"`
}

type GenerateCommand struct {
	Privkey string `long:"priv" description:"Path to save privkey" required:"true"`
	Pubkey string `long:"pub" description:"Path to save pubkey if you want" required:"true"`
	Sshkey string `long:"ssh" description:"Path to save pubkey in SSH format if you want" required:"false"`
}

type RunCommand struct {
	Host string `short:"H" long:"host" description:"Desired host" required:"false"`
	Port string `short:"P" long:"port" description:"Local tunnel port" required:"true"`
}

func (x *GenerateCommand) Execute(args []string) error {
	generate(x.Privkey, x.Pubkey, x.Sshkey)
	return nil
}

func (x *TunnelCommand) Execute(args []string) error {
	tunnel(x.Host, x.LocalHost, x.LocalPort, x.Key, x.ServerPort, x.ServerIp)
	return nil
}

func (x *RunCommand) Execute(args []string) error {
	var err error

	fmt.Print("Input server address: ")
	_, err = fmt.Scanf("%s\n", &Configuration.ServerAddr)
	checkError(err)


	fmt.Print("Input server port: ")
	_, err = fmt.Scanf("%d\n", &Configuration.ServerPort)
	checkError(err)

	fmt.Print("Input access password: ")
	_, err = fmt.Scanf("%s\n", &Configuration.Password)
	checkError(err)

	configJSON, _ := json.Marshal(Configuration)
	err = ioutil.WriteFile("config.json", configJSON, 0644)
	checkError(err)

	return nil
}