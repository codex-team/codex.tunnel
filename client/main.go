package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/function61/gokit/bidipipe"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strings"
)

func main() {
	//generate("key")
	tunnel()
}

func generate(name string) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	savePublicPEMKey(fmt.Sprintf("%s.pub", name), publicKey)
	saveSSHKey(fmt.Sprintf("%s.public.pem", name), &publicKey)
	savePEMKey(fmt.Sprintf("%s.private.pem", name), key)
}

func tunnel() {
	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: 9000,
	}

	serverEndpoint := &Endpoint{
		Host: "localhost",
		Port: 8022,
	}

	sshConfig := &ssh.ClientConfig{
		User: "codex",
		Auth: []ssh.AuthMethod{
			privateKeyFile("../tun_key.pem"),
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
	PORT := strings.Split(remoteAddr.String(), ":")[1]

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

	session.Start("nostr")
	in.Write([]byte("myhost\n"))
	in.Write([]byte(PORT + "\n"))

	//var (
	//	line string
	//	r    = bufio.NewReader(out)
	//	output []byte
	//)
	//for {
	//	b, err := r.ReadByte()
	//	if err != nil {
	//		break
	//	}
	//	output = append(output, b)
	//	if b == byte('\n') {
	//		line = ""
	//		continue
	//	}
	//	line += string(b)
	//	log.Printf(line)
	//	if line == "lol" {
	//		log.Printf("send")
	//		in.Write([]byte("myhost\n"))
	//		in.Write([]byte(remoteAddr.String() + "\n"))
	//		break
	//	}
	//}

	defer session.Close()

	go func() {
		defer listener.Close()

		log.Println(fmt.Sprintf("listening remote %s", remoteAddr))

		// handle incoming connections on reverse forwarded tunnel
		for {
			client, err := listener.Accept()
			if err != nil {
				return
			}

			go handleClient2(client, localEndpoint.String())
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

func handleClient2(client net.Conn, localEndpoint string) {
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