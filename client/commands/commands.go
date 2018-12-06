package commands

import (
	"fmt"
	"os"
)

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

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}