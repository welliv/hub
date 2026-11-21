package config

// DefaultElectrumServer is the default ELECTRUM_SERVER (a mainnet server).
const DefaultElectrumServer = "ssl://electrum.getalbypro.com:50002"

// DefaultMempoolApi returns the mempool.space API URL for a network, used when
// MEMPOOL_API is not set. Block explorer links are derived from it (see
// GetMempoolUrl), so on testnet they point to mempool.space/testnet.
// Networks without a public mempool.space instance (e.g. regtest) keep the
// mainnet default.
func DefaultMempoolApi(network string) string {
	switch network {
	case "testnet", "testnet3":
		return "https://mempool.space/testnet/api"
	case "testnet4":
		return "https://mempool.space/testnet4/api"
	case "signet":
		return "https://mempool.space/signet/api"
	default:
		return "https://mempool.space/api"
	}
}

// IsMainnet reports whether a network name refers to bitcoin mainnet.
func IsMainnet(network string) bool {
	return network == "bitcoin" || network == "mainnet"
}

// applyNetworkDefaults fills in network-dependent defaults. An explicitly set
// MEMPOOL_API always wins.
func (c *AppConfig) applyNetworkDefaults() {
	if c.MempoolApi == "" {
		c.MempoolApi = DefaultMempoolApi(c.Network)
	}
}
