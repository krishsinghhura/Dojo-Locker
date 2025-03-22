package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GeneratePandPkey(oldhash string) (string, string, error) {
	pvtKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) //elliptic.P256 is an algorithm through which we will firstly generate a random hash
	if err != nil {
		fmt.Println("Error generating key:", err)
		return "", "", err
	}

	pvtBytes, err := x509.MarshalECPrivateKey(pvtKey) //MarshalECPrivateKey it encodes the key in x509 format, which is further use for saftey
	//we cant expose the private key openly thats why we need to encode it in x509 format
	if err != nil {
		fmt.Println("Error marshalling private key:", err)
		return "", "", err
	}

	pvtPem := pem.EncodeToMemory(&pem.Block{ //PEM (Privacy-Enhanced Mail)
		//making a text file to store the pvtkey but with embedding the hash of the document
		Type:  "EC PRIVATE KEY",
		Bytes: pvtBytes,
		Headers: map[string]string{
			"Hash": oldhash,
		},
	})

	//now for public
	pbcBytes, err := x509.MarshalPKIXPublicKey(&pvtKey.PublicKey) //converting the public key into x509 same as we did in pvt key
	if err != nil {
		fmt.Println("Error marshalling public key:", err)
		return "", "", err
	}

	pbcPem := pem.EncodeToMemory(&pem.Block{
		//saves it in a file but with embedding the hash of the document
		Type:  "PUBLIC KEY",
		Bytes: pbcBytes,
		Headers: map[string]string{
			"Hash": oldhash,
		},
	})

	return string(pvtPem), string(pbcPem), nil
}
