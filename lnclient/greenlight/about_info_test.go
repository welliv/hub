package greenlight

import "testing"

func TestGetLiquiditySourceLsps2PaymentSizeRangeMsat(t *testing.T) {
	svc, cleanup := newTestService(t, newMockNode())
	defer cleanup()

	min, max := svc.GetLiquiditySourceLsps2PaymentSizeRangeMsat()
	if min != nil || max != nil {
		t.Fatalf("unconfigured bounds must be nil, got min=%v max=%v", min, max)
	}

	svc.config.JitMinPaymentSizeMsat = 2_501_000
	svc.config.JitMaxPaymentSizeMsat = 16_000_000_000
	min, max = svc.GetLiquiditySourceLsps2PaymentSizeRangeMsat()
	if min == nil || *min != 2_501_000 {
		t.Fatalf("min: got %v, want 2501000", min)
	}
	if max == nil || *max != 16_000_000_000 {
		t.Fatalf("max: got %v, want 16000000000", max)
	}
}
