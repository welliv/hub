package api

import (
	"archive/zip"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func greenlightCertZip(t *testing.T, extraFile bool) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}, &x509.Certificate{SerialNumber: big.NewInt(1)}, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	crt := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	priv := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range map[string][]byte{"client.crt": crt, "client-key.pem": priv} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(content); err != nil {
			t.Fatal(err)
		}
	}
	if extraFile {
		w, err := zw.Create("extra.txt")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("extra"))
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestParseGreenlightDeveloperCertZip(t *testing.T) {
	crt, key, err := parseGreenlightDeveloperCertZip(greenlightCertZip(t, false))
	if err != nil {
		t.Fatalf("valid zip rejected: %v", err)
	}
	if crt == "" || key == "" {
		t.Fatal("expected PEM contents")
	}
	if _, _, err := parseGreenlightDeveloperCertZip(greenlightCertZip(t, true)); err == nil {
		t.Fatal("extra file must be rejected")
	}
}
