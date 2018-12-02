package main

import (
	"fmt"
	"github.com/function61/gokit/bidipipe"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func handleClient(client net.Conn, localEndpoint string) {
	defer client.Close()

	log.Println(fmt.Sprintf("%s connected", client.RemoteAddr()))
	defer log.Println("closed")

	remote, err := net.Dial("tcp", localEndpoint)
	if err != nil {
		log.Println(fmt.Sprintf("dial INTO local service error: %s", err.Error()))
		return
	}

	if err := bidipipe.Pipe(client, "client", remote, "remote"); err != nil {
		log.Println(err.Error())
	}
}

func privateKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH private key file %s", file))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH private key file %s", file))
		return nil
	}
	return ssh.PublicKeys(key)
}