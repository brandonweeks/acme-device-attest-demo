package main

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"os/exec"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/go-attestation/attest"
	x509ext "github.com/google/go-attestation/x509"
	"github.com/google/go-tpm-tools/simulator"
	"golang.org/x/crypto/acme"
)

const (
	accountKeyFile = "account.key"
)

var (
	caAddress    = flag.String("ca_address", "https://ca.attestation.dev/acme/acme/directory", "URL of ACME directory endpoint")
	serialNumber = flag.String("serial", "12345", "Device serial number")
	useSimulator = flag.Bool("sim", false, "Use a simulated TPM")
)

func accountKey() (crypto.Signer, error) {
	if _, err := os.Stat(accountKeyFile); err == nil {
		der, err := os.ReadFile(accountKeyFile)
		if err != nil {
			return nil, err
		}
		key, err := x509.ParsePKCS8PrivateKey(der)
		if err != nil {
			return nil, err
		}
		switch t := key.(type) {
		case *rsa.PrivateKey:
			return t, nil
		case *ecdsa.PrivateKey:
			return t, nil
		default:
			return nil, fmt.Errorf("unsupported private key type: %T", key)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		der, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(accountKeyFile, der, 0600); err != nil {
			return nil, err
		}
		return key, nil
	} else {
		return nil, err
	}
}

func akCert(ak *attest.AK) ([]byte, error) {
	akRootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	akRootTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
	}
	permID := x509ext.PermanentIdentifier{
		IdentifierValue: *serialNumber,
		Assigner:        asn1.ObjectIdentifier{0, 1, 2, 3, 4},
	}
	san := &x509ext.SubjectAltName{
		PermanentIdentifiers: []x509ext.PermanentIdentifier{
			permID,
		},
	}
	ext, err := x509ext.MarshalSubjectAltName(san)
	if err != nil {
		return nil, err
	}
	akTemplate := &x509.Certificate{
		SerialNumber:    big.NewInt(2),
		ExtraExtensions: []pkix.Extension{ext},
	}
	akPub, err := attest.ParseAKPublic(attest.TPMVersion20, ak.AttestationParameters().Public)
	if err != nil {
		return nil, err
	}
	akCert, err := x509.CreateCertificate(rand.Reader, akTemplate, akRootTemplate, akPub.Public, akRootKey)
	if err != nil {
		return nil, err
	}
	return akCert, nil
}

type AttestationObject struct {
	Format       string                 `json:"fmt"`
	AttStatement map[string]interface{} `json:"attStmt,omitempty"`
}

func attestationStatement(key *attest.Key, akCert []byte) ([]byte, error) {
	params := key.CertificationParameters()

	obj := &AttestationObject{
		Format: "tpm",
		AttStatement: map[string]interface{}{
			"ver":      "2.0",
			"alg":      int64(-257), // AlgRS256
			"x5c":      []interface{}{akCert},
			"sig":      params.CreateSignature,
			"certInfo": params.CreateAttestation,
			"pubArea":  params.Public,
		},
	}
	b, err := cbor.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func csr(key *attest.Key) ([]byte, error) {
	permID := x509ext.PermanentIdentifier{
		IdentifierValue: *serialNumber,
		Assigner:        asn1.ObjectIdentifier{0, 1, 2, 3, 4},
	}
	san := &x509ext.SubjectAltName{
		PermanentIdentifiers: []x509ext.PermanentIdentifier{
			permID,
		},
	}
	ext, err := x509ext.MarshalSubjectAltName(san)
	if err != nil {
		return nil, err
	}
	tmpl := &x509.CertificateRequest{
		ExtraExtensions: []pkix.Extension{ext},
	}
	privKey, err := key.Private(key.Public())
	if err != nil {
		return nil, err
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, tmpl, privKey.(crypto.Signer))
	if err != nil {
		return nil, err
	}
	return der, nil
}

type simulatorChannel struct {
	io.ReadWriteCloser
}

func (simulatorChannel) MeasurementLog() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func tpmInit() (*attest.Key, []byte, error) {
	config := &attest.OpenConfig{}
	if *useSimulator {
		sim, err := simulator.Get()

		if err != nil {
			return nil, nil, err
		}
		config.CommandChannel = simulatorChannel{sim}
	}
	tpm, err := attest.OpenTPM(config)
	if err != nil {
		return nil, nil, err
	}
	ak, err := tpm.NewAK(nil)
	if err != nil {
		return nil, nil, err
	}
	key, err := tpm.NewKey(ak, nil)
	if err != nil {
		return nil, nil, err
	}
	akCert, err := akCert(ak)
	if err != nil {
		return nil, nil, err
	}
	return key, akCert, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Cloud Shell: hack to give the unprivileged user access to the TPM
	// resource manager.
	if os.Getenv("CLOUD_SHELL") == "true" {
		cmd := exec.Command("sudo", "chmod", "777", "/dev/tpmrm0")
		if b, err := cmd.CombinedOutput(); err != nil {
			log.Fatal(string(b))
		}
	}

	accountKey, err := accountKey()
	if err != nil {
		log.Fatal(err)
	}

	client := acme.Client{
		Key:          accountKey,
		DirectoryURL: *caAddress,
	}

	account := &acme.Account{}
	_, err = client.Register(ctx, account, func(tosURL string) bool { panic("") })
	if err != nil && !errors.Is(err, acme.ErrAccountAlreadyExists) {
		log.Fatal(err)
	}

	id := []acme.AuthzID{
		{
			Type:  "permanent-identifier",
			Value: *serialNumber,
		},
	}
	opts := []acme.OrderOption{}
	order, err := client.AuthorizeOrder(ctx, id, opts...)
	if err != nil {
		log.Fatal(err)
	}

	certKey, akCert, err := tpmInit()
	if err != nil {
		log.Fatal(err)
	}

	for _, authzURL := range order.AuthzURLs {
		authz, err := client.GetAuthorization(ctx, authzURL)
		if err != nil {
			log.Fatal(err)
		}

		for _, chal := range authz.Challenges {
			payload, err := attestationStatement(certKey, akCert)
			if err != nil {
				log.Fatal(err)
			}
			req := struct {
				AttStmt []byte `json:"attStmt"`
			}{
				payload,
			}
			if _, err = client.AcceptWithPayload(ctx, chal, req); err != nil {
				log.Fatal(err)
			}
		}
	}

	csr, err := csr(certKey)
	if err != nil {
		log.Fatal(err)
	}
	der, _, err := client.CreateOrderCert(ctx, order.FinalizeURL, csr, false)
	if err != nil {
		log.Fatal(err)
	}
	b := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der[0],
	}
	fmt.Println(string(pem.EncodeToMemory(b)))
}
