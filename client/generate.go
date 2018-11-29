package main

import (
"crypto/rsa"
"crypto/x509"
"encoding/asn1"
"encoding/pem"
"fmt"
"golang.org/x/crypto/ssh"
"io/ioutil"
"os"
)

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
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

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}