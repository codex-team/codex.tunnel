package main

var opts struct {
	Generate GenerateCommand `command:"generate" description:"Generate rsa keys"`
	Tunnel TunnelCommand `command:"tunnel" description:"Start tunnel"`
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

var GenerateCfg GenerateCommand
var TunnelCfg TunnelCommand

func (x *GenerateCommand) Execute(args []string) error {
	generate(x.Privkey, x.Pubkey, x.Sshkey)
	return nil
}

func (x *TunnelCommand) Execute(args []string) error {
	tunnel(x.Host, x.LocalHost, x.LocalPort, x.Key, x.ServerPort, x.ServerIp)
	return nil
}