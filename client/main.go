package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"github.com/jessevdk/go-flags"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
)

func main() {
	_, err := flags.Parse(&opts)

	if err != nil {
		//log.Fatalln(err)
		os.Exit(0)
	}
}

func generate(privkeyName string, pubkeyName string, sshkeyName string) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	savePublicPEMKey(fmt.Sprintf("%s", pubkeyName), publicKey)
	savePEMKey(fmt.Sprintf("%s", privkeyName), key)

	if sshkeyName != "" {
		saveSSHKey(fmt.Sprintf("%s", sshkeyName), &publicKey)
	}
}

func tunnel(ServerHost string, LocalHost string, LocalPort int, privateKeyFilename string) {
	localEndpoint := &Endpoint{
		Host: LocalHost,
		Port: LocalPort,
	}

	serverEndpoint := &Endpoint{
		Host: "tun.ifmo.su",
		Port: 17022,
	}

	sshConfig := &ssh.ClientConfig{
		User: "codex",
		Auth: []ssh.AuthMethod{
			privateKeyFile(privateKeyFilename),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ctx, _ := context.WithCancel(context.Background())

	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
	}

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
	in.Write([]byte(ServerHost + "\n"))
	in.Write([]byte(ServerPort + "\n"))

	defer session.Close()

	go func() {
		defer listener.Close()

		log.Println(fmt.Sprintf("tunneling %s.tun.ifmo.su --> %s", ServerHost, localEndpoint.String()))

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
			return
		default:
		}
	}

	serverConn.Close()
}