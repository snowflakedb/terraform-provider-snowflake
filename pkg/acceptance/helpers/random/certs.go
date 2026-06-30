package random

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/youmark/pkcs8"
)

// GenerateX509 returns base64 encoded certificate on a single line without the leading -----BEGIN CERTIFICATE----- and ending -----END CERTIFICATE----- markers.
func GenerateX509(t *testing.T) string {
	t.Helper()
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Company, INC."},
		},
		NotAfter:    time.Now().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	require.NoError(t, err)

	return encode(t, "CERTIFICATE", caBytes)
}

// GenerateRSAPublicKey returns an RSA public key without BEGIN and END markers, and key's hash.
func GenerateRSAPublicKey(t *testing.T) (string, string) {
	t.Helper()
	key := GenerateRSAPrivateKey(t)

	return GenerateRSAPublicKeyFromPrivateKey(t, key)
}

// GenerateRSAPublicKeyFromPrivateKey returns an RSA public key without BEGIN and END markers, and key's hash.
func GenerateRSAPublicKeyFromPrivateKey(t *testing.T, key *rsa.PrivateKey) (string, string) {
	t.Helper()

	pub := key.Public()
	b, err := x509.MarshalPKIXPublicKey(pub.(*rsa.PublicKey))
	require.NoError(t, err)
	return encode(t, "RSA PUBLIC KEY", b), hash(t, b)
}

// GenerateRSAPrivateKey returns an RSA private key.
func GenerateRSAPrivateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key
}

// GenerateRSAKeyPair returns an RSA private key (unencrypted and encrypted), RSA public key without BEGIN and END markers, and key's hash.
func GenerateRSAKeyPair(t *testing.T, pass string) (string, string, string, string) {
	t.Helper()

	privateKey := GenerateRSAPrivateKey(t)
	unencryptedDer, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)
	privBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: unencryptedDer,
	}
	unencrypted := string(pem.EncodeToMemory(&privBlock))
	encrypted := encrypt(t, privateKey, pass)

	publicKey, keyHash := GenerateRSAPublicKeyFromPrivateKey(t, privateKey)
	return unencrypted, encrypted, publicKey, keyHash
}

// GenerateRSAPrivateKeyEncrypted returns a PEM-encoded pair of unencrypted and encrypted key with a given password
func GenerateRSAPrivateKeyEncrypted(t *testing.T, password string) (unencrypted, encrypted string) {
	t.Helper()
	rsaPrivateKey := GenerateRSAPrivateKey(t)
	unencryptedDer, err := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)
	require.NoError(t, err)
	privBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: unencryptedDer,
	}
	unencrypted = string(pem.EncodeToMemory(&privBlock))
	encrypted = encrypt(t, rsaPrivateKey, password)

	return
}

func encrypt(t *testing.T, rsaPrivateKey *rsa.PrivateKey, pass string) string {
	t.Helper()

	encryptedDer, err := pkcs8.MarshalPrivateKey(rsaPrivateKey, []byte(pass), &pkcs8.Opts{
		Cipher: pkcs8.AES256CBC,
		KDFOpts: pkcs8.PBKDF2Opts{
			SaltSize: 16, IterationCount: 2000, HMACHash: crypto.SHA256,
		},
	})
	require.NoError(t, err)
	privEncryptedBlock := pem.Block{
		Type:  "ENCRYPTED PRIVATE KEY",
		Bytes: encryptedDer,
	}
	return string(pem.EncodeToMemory(&privEncryptedBlock))
}

func hash(t *testing.T, b []byte) string {
	t.Helper()
	hash := sha256.Sum256(b)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func encode(t *testing.T, pemType string, b []byte) string {
	t.Helper()
	buffer := new(bytes.Buffer)
	err := pem.Encode(
		buffer,
		&pem.Block{
			Type:  pemType,
			Bytes: b,
		},
	)
	require.NoError(t, err)
	cert := strings.TrimPrefix(buffer.String(), fmt.Sprintf("-----BEGIN %s-----\n", pemType))
	cert = strings.TrimSuffix(cert, fmt.Sprintf("-----END %s-----\n", pemType))
	return cert
}
