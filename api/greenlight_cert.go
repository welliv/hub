package api

import (
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	maxGreenlightCertZipSize = 128 << 10
	maxGreenlightCertPEMSize = 32 << 10
)

type GreenlightCertStatus struct {
	Stage    string `json:"stage"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
}

func (api *api) setGreenlightCertStatus(stage string, progress int, message string, err error) {
	api.greenlightCertStatusMu.Lock()
	defer api.greenlightCertStatusMu.Unlock()
	api.greenlightCertStatus = GreenlightCertStatus{
		Stage:    stage,
		Progress: progress,
		Message:  message,
	}
	if err != nil {
		api.greenlightCertStatus.Stage = "failed"
		api.greenlightCertStatus.Error = err.Error()
	}
}

func (api *api) GetGreenlightCertStatus() GreenlightCertStatus {
	api.greenlightCertStatusMu.Lock()
	status := api.greenlightCertStatus
	api.greenlightCertStatusMu.Unlock()

	if status.Stage == "" || status.Stage == "idle" {
		uploaded, err := api.cfg.Get("GreenlightCertUploaded", "")
		if err == nil && uploaded == "true" {
			return GreenlightCertStatus{Stage: "ready", Progress: 100, Message: "Certificate stored"}
		}
		return GreenlightCertStatus{Stage: "idle", Progress: 0, Message: "No certificate uploaded"}
	}
	return status
}

// StartGreenlightCertUpload accepts a zip with client.crt and client-key.pem,
// validates it in the background, and stores both PEMs encrypted with the
// unlock password. Progress is exposed through GetGreenlightCertStatus.
func (api *api) StartGreenlightCertUpload(unlockPassword string, r io.Reader) error {
	if unlockPassword == "" {
		return errors.New("unlock password is required")
	}
	current := api.GetGreenlightCertStatus()
	if current.Stage == "uploading" || current.Stage == "validating" || current.Stage == "storing" {
		return errors.New("certificate upload already in progress")
	}

	api.setGreenlightCertStatus("uploading", 10, "Uploading certificate", nil)
	raw, err := io.ReadAll(io.LimitReader(r, maxGreenlightCertZipSize+1))
	if err != nil {
		api.setGreenlightCertStatus("failed", 0, "Upload failed", err)
		return err
	}
	if len(raw) > maxGreenlightCertZipSize {
		err = errors.New("certificate zip is too large")
		api.setGreenlightCertStatus("failed", 0, "Upload failed", err)
		return err
	}

	go func() {
		api.setGreenlightCertStatus("validating", 35, "Validating certificate", nil)
		crt, key, err := parseGreenlightDeveloperCertZip(raw)
		if err != nil {
			api.setGreenlightCertStatus("failed", 0, "Invalid certificate", err)
			return
		}

		api.setGreenlightCertStatus("storing", 70, "Storing encrypted certificate", nil)
		if err := api.cfg.SetUpdate("GreenlightNobodyCrt", crt, unlockPassword); err != nil {
			api.setGreenlightCertStatus("failed", 0, "Could not store certificate", err)
			return
		}
		if err := api.cfg.SetUpdate("GreenlightNobodyKey", key, unlockPassword); err != nil {
			api.setGreenlightCertStatus("failed", 0, "Could not store certificate key", err)
			return
		}
		if err := api.cfg.SetUpdate("GreenlightCertUploaded", "true", ""); err != nil {
			api.setGreenlightCertStatus("failed", 0, "Could not mark certificate ready", err)
			return
		}
		api.setGreenlightCertStatus("ready", 100, "Certificate ready", nil)
	}()

	return nil
}

func parseGreenlightDeveloperCertZip(raw []byte) (string, string, error) {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return "", "", fmt.Errorf("invalid zip: %w", err)
	}
	if len(zr.File) != 2 {
		return "", "", fmt.Errorf("expected exactly 2 files in zip, got %d", len(zr.File))
	}

	files := map[string]string{}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() || f.Mode()&os.ModeSymlink != 0 {
			return "", "", fmt.Errorf("unexpected file entry %q", f.Name)
		}
		if f.Name != "client.crt" && f.Name != "client-key.pem" {
			return "", "", fmt.Errorf("unexpected file %q: expected client.crt and client-key.pem", f.Name)
		}
		rc, err := f.Open()
		if err != nil {
			return "", "", err
		}
		data, err := io.ReadAll(io.LimitReader(rc, maxGreenlightCertPEMSize+1))
		_ = rc.Close()
		if err != nil {
			return "", "", err
		}
		if len(data) > maxGreenlightCertPEMSize {
			return "", "", fmt.Errorf("%s is too large", f.Name)
		}
		files[f.Name] = string(data)
	}

	crtPEM, ok := files["client.crt"]
	if !ok {
		return "", "", errors.New("missing client.crt")
	}
	keyPEM, ok := files["client-key.pem"]
	if !ok {
		return "", "", errors.New("missing client-key.pem")
	}

	cert, err := parseGreenlightCertPEM(crtPEM)
	if err != nil {
		return "", "", fmt.Errorf("client.crt: %w", err)
	}
	if err := validateGreenlightKeyPEM(keyPEM, cert); err != nil {
		return "", "", fmt.Errorf("client-key.pem: %w", err)
	}
	return crtPEM, keyPEM, nil
}

func parseGreenlightCertPEM(value string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(strings.TrimSpace(value)))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("not a PEM certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

func validateGreenlightKeyPEM(value string, cert *x509.Certificate) error {
	block, _ := pem.Decode([]byte(strings.TrimSpace(value)))
	if block == nil || !strings.Contains(block.Type, "PRIVATE KEY") {
		return errors.New("not a PEM private key")
	}

	var signer crypto.Signer
	switch block.Type {
	case "PRIVATE KEY":
		parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		var ok bool
		signer, ok = parsed.(crypto.Signer)
		if !ok {
			return errors.New("unsupported private key type")
		}
	case "EC PRIVATE KEY":
		parsed, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		signer = parsed
	case "RSA PRIVATE KEY":
		parsed, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		signer = parsed
	default:
		return fmt.Errorf("unsupported private key type %q", block.Type)
	}

	pub, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		return err
	}
	if !bytes.Equal(pub, cert.RawSubjectPublicKeyInfo) {
		return errors.New("private key does not match certificate")
	}
	return nil
}
