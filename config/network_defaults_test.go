package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultMempoolApi(t *testing.T) {
	for network, expected := range map[string]string{
		"bitcoin":  "https://mempool.space/api",
		"mainnet":  "https://mempool.space/api",
		"testnet":  "https://mempool.space/testnet/api",
		"testnet3": "https://mempool.space/testnet/api",
		"testnet4": "https://mempool.space/testnet4/api",
		"signet":   "https://mempool.space/signet/api",
		"regtest":  "https://mempool.space/api",
		"":         "https://mempool.space/api",
	} {
		assert.Equal(t, expected, DefaultMempoolApi(network), network)
	}
}

func TestApplyNetworkDefaults(t *testing.T) {
	t.Run("unset mempool api follows the network", func(t *testing.T) {
		env := &AppConfig{Network: "testnet"}
		env.applyNetworkDefaults()
		assert.Equal(t, "https://mempool.space/testnet/api", env.MempoolApi)
	})
	t.Run("explicit mempool api wins", func(t *testing.T) {
		env := &AppConfig{Network: "testnet", MempoolApi: "https://mempool.example.com/api"}
		env.applyNetworkDefaults()
		assert.Equal(t, "https://mempool.example.com/api", env.MempoolApi)
	})
}

func TestGetMempoolUrl_FollowsNetwork(t *testing.T) {
	for network, expected := range map[string]string{
		"bitcoin": "https://mempool.space",
		"testnet": "https://mempool.space/testnet",
		"signet":  "https://mempool.space/signet",
	} {
		env := &AppConfig{Network: network}
		env.applyNetworkDefaults()
		cfg := &config{Env: env}
		assert.Equal(t, expected, cfg.GetMempoolUrl(), network)
	}
	// also when the env was never passed through init
	cfg := &config{Env: &AppConfig{Network: "testnet"}}
	assert.Equal(t, "https://mempool.space/testnet", cfg.GetMempoolUrl())
}

func TestDefaultElectrumServerMatchesEnvDefault(t *testing.T) {
	field, ok := reflect.TypeOf(AppConfig{}).FieldByName("ElectrumServer")
	assert.True(t, ok)
	assert.Equal(t, DefaultElectrumServer, field.Tag.Get("default"))
}
