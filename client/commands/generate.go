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

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func (x *GenerateCommand) Execute(args []string) error {
	return generateKeys(x.Privkey, x.Pubkey, x.Sshkey)
}

func generateKeys(Privkey string, Pubkey string, Sshkey string) error {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return err
	}

	publicKey := key.PublicKey

	err = savePublicPEMKey(fmt.Sprintf("%s", Pubkey), publicKey)
	if err != nil {
		return err
	}
	err = savePEMKey(fmt.Sprintf("%s", Privkey), key)
	if err != nil {
		return err
	}

	if Sshkey != "" {
		err = saveSSHKey(fmt.Sprintf("%s", Sshkey), &publicKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func savePEMKey(fileName string, key *rsa.PrivateKey) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return pem.Encode(outFile, privateKey)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) error {
	asn1Bytes, err := asn1.Marshal(pubkey)
	if err != nil {
		return err
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer pemfile.Close()

	return pem.Encode(pemfile, pemkey)
}

func saveSSHKey(fileName string, key *rsa.PublicKey) error {
	publicRsaKey, err := ssh.NewPublicKey(key)
	if err != nil {
		return err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	return ioutil.WriteFile(fileName, pubKeyBytes, 0600)
}

func RandASCIIBytes(n int) []byte {
	output := make([]byte, n)
	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)
	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letterBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])
		// random % 64
		randomPos := random % uint8(l)
		// put into output
		output[pos] = letterBytes[randomPos]
	}
	return output
}