package greenlight

// About-page metadata for the Greenlight backend.
//
// The Hub's /api/info and /api/jit-channels-info endpoints ask the active
// backend for these values through small optional interfaces. Keeping the
// answers here means the About page can show the same shape for every backend
// without hardcoding a backend type in the API layer.
// The gl-plugin selects an LSP from connected peers at invoice time and
// does not expose which one it chose. Say that plainly rather than inventing
// a pubkey or reusing an LDK catalog entry from another network.
const (
	greenlightLiquiditySource = "Greenlight plugin-selected LSP"
)

// NOTE: the chain data source is intentionally NOT provided here.
// It is a Hub-level setting (BITCOIND_RPC_HOST / ESPLORA_SERVER /
// ELECTRUM_SERVER) shared by every backend, so api.GetInfo derives it from
// config for Greenlight exactly as it does for LDK. Reporting "hosted" here
// would hide the real electrum/esplora endpoint from the About page.

// GetLiquiditySourceLsps2 implements the api.lsps2SourceProvider capability.
func (g *GreenlightService) GetLiquiditySourceLsps2() string {
	return greenlightLiquiditySource
}

// GetLiquiditySourceLsps2PaymentSizeRangeMsat implements the
// api.lsps2PaymentSizeRangeProvider capability.
//
// No bounds are published for the hosted LSP, so return nil for both: the
// frontend omits the "JIT payment size" line rather than fabricating a range
// from a different network. Once a real quote is observed it can be cached
// here and reported.
func (g *GreenlightService) GetLiquiditySourceLsps2PaymentSizeRangeMsat() (*uint64, *uint64) {
	return nil, nil
}
