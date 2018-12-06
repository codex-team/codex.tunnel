package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

func (x *GenerateCommand) Execute(args []string) error {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	savePublicPEMKey(fmt.Sprintf("%s", x.Pubkey), publicKey)
	savePEMKey(fmt.Sprintf("%s", x.Privkey), key)

	if x.Sshkey != "" {
		saveSSHKey(fmt.Sprintf("%s", x.Sshkey), &publicKey)
	}

	return nil
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func saveSSHKey(fileName string, key *rsa.PublicKey) {
	publicRsaKey, err := ssh.NewPublicKey(key)
	checkError(err)
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	err = ioutil.WriteFile(fileName, pubKeyBytes, 0600)
	checkError(err)
}

