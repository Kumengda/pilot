package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func NewCert() error {
	// 生成一个 ECDSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 构造证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"pilotProxy"},
			OrganizationalUnit: []string{"Service Infrastructure Department"},
			Locality:           []string{"Chengdu"},
			Province:           nil,
			StreetAddress:      []string{"Chengdu"},
			PostalCode:         nil,
			SerialNumber:       "",
			CommonName:         "Root CA For Pilot",
			Names:              nil,
			ExtraNames:         nil,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(100, 0, 0),
		KeyUsage:    0,
		ExtKeyUsage: nil,
	}

	// 使用证书模板和公钥生成证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// 将证书和私钥写入文件
	certOut, err := os.Create("pilot.crt")
	if err != nil {
		return err
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}
	err = certOut.Close()
	if err != nil {
		return err
	}

	keyOut, err := os.Create("pilot.key")
	if err != nil {
		return err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err != nil {
		return err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		return err
	}
	err = keyOut.Close()
	if err != nil {
		return err
	}
	return nil
}
