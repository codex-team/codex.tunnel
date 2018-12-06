package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/function61/gokit/bidipipe"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"strings"
)

type Endpoint struct {
	Host string
	Port int
}

const DefaultUser = "codex"

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func make_tunnel(Host string, LocalHost string, LocalPort int, ServerHost string, ServerOutPort int, Key string) error {
	localEndpoint := &Endpoint{
		Host: LocalHost,
		Port: LocalPort,
	}

	u, err := url.Parse(ServerHost)
	if err != nil {
		return err
	}

	ServerHostOnly, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return err
	}

	serverEndpoint := &Endpoint{
		Host: ServerHostOnly,
		Port: ServerOutPort,
	}

	sshConfig := &ssh.ClientConfig{
		User: DefaultUser,
		Auth: []ssh.AuthMethod{
			privateKeyFile(Key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ctx, _ := context.WithCancel(context.Background())

	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
	}

	defer serverConn.Close()

	listener, err := serverConn.Listen("tcp", fmt.Sprintf("%s:0", serverEndpoint.Host))
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
	}
	remoteAddr := listener.Addr()
	ServerPort := strings.Split(remoteAddr.String(), ":")[1]

	session, err := serverConn.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	in, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	if strings.Contains(stdoutBuf.String(), "[error]") {
		log.Fatalln(fmt.Sprintf("%s", stdoutBuf.String()))
	}

	session.Start("auth")
	in.Write([]byte(Host + "\n"))
	in.Write([]byte(ServerPort + "\n"))

	defer session.Close()

	go func() {
		defer listener.Close()

		log.Println(fmt.Sprintf("tunneling %s.%s --> %s", Host, ServerHostOnly, localEndpoint.String()))

		for {
			client, err := listener.Accept()
			if err != nil {
				return
			}

			go handleClient(client, localEndpoint.String())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	return nil
}

func (x *TunnelCommand) Execute(args []string) error {
	return make_tunnel(x.Host, x.LocalHost, x.LocalPort, x.ServerIp, x.ServerPort, x.Key)
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