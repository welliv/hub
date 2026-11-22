package greenlight

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"

	"github.com/getAlby/hub/lnclient"
)

// validIncomingPreimage checks a hosted node's payment event against the
// cryptographic invariant of Lightning: SHA-256(preimage) == payment_hash.
//
// A missing preimage is not treated as a lie — some node versions omit it on
// reads — but an explicitly present, mismatching preimage must never be
// allowed to credit a wallet.
func validIncomingPreimage(tx *lnclient.Transaction) bool {
	if tx == nil {
		return false
	}
	if tx.PaymentHash == "" || tx.Preimage == "" {
		return true
	}
	hash, err := hex.DecodeString(tx.PaymentHash)
	if err != nil || len(hash) != sha256.Size {
		return false
	}
	preimage, err := hex.DecodeString(tx.Preimage)
	if err != nil || len(preimage) != sha256.Size {
		return false
	}
	sum := sha256.Sum256(preimage)
	return subtle.ConstantTimeCompare(sum[:], hash) == 1
}
