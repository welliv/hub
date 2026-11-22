package greenlight

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/getAlby/hub/lnclient"
)

func TestValidIncomingPreimage(t *testing.T) {
	preimage := make([]byte, 32)
	copy(preimage, []byte("test preimage"))
	hash := sha256.Sum256(preimage)

	tx := &lnclient.Transaction{
		PaymentHash: hex.EncodeToString(hash[:]),
		Preimage:    hex.EncodeToString(preimage),
	}
	if !validIncomingPreimage(tx) {
		t.Fatal("valid preimage rejected")
	}

	bad := &lnclient.Transaction{
		PaymentHash: hex.EncodeToString(hash[:]),
		Preimage:    hex.EncodeToString(make([]byte, 32)),
	}
	if validIncomingPreimage(bad) {
		t.Fatal("mismatched preimage accepted")
	}

	missing := &lnclient.Transaction{PaymentHash: hex.EncodeToString(hash[:])}
	if !validIncomingPreimage(missing) {
		t.Fatal("missing preimage must not be treated as a lie")
	}
}
