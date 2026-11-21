package ldk

import (
	"testing"

	"github.com/getAlby/ldk-node-go/ldk_node"
	"github.com/stretchr/testify/assert"
)

func TestLdkNetwork(t *testing.T) {
	assert.Equal(t, ldk_node.NetworkBitcoin, ldkNetwork("bitcoin"))
	assert.Equal(t, ldk_node.NetworkTestnet, ldkNetwork("testnet"))
	assert.Equal(t, ldk_node.NetworkTestnet, ldkNetwork("testnet3"))
	assert.Equal(t, ldk_node.NetworkSignet, ldkNetwork("signet"))
	assert.Equal(t, ldk_node.NetworkRegtest, ldkNetwork("regtest"))
}
