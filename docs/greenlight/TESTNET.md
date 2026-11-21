# Testnet trial (testnet3)

Greenlight hosts two networks: `testnet` (testnet3) and `bitcoin`. Use testnet to try the build without
real funds. Testnet coins have no value.

## Run the hub on testnet

Build as in [mainnet-trial/README.md](mainnet-trial/README.md), then start the hub with:

```bash
export NETWORK=testnet LN_BACKEND_TYPE=GREENLIGHT
export GL_NOBODY_CRT=/path/to/devcert/client.crt GL_NOBODY_KEY=/path/to/devcert/client-key.pem
export GREENLIGHT_GLCLI_PATH=$(command -v glcli)      # glcli 0.1.3, see mainnet-trial/README.md
export WORK_DIR=$HOME/.albyhub-gl-testnet PORT=127.0.0.1:8080
./albyhub
```

Open the UI, choose a password, and the hub registers a new testnet node under your developer
certificate (about 10 s). Keep the developer certificate out of chats and repositories.

To compare with the default backend, run a second hub with its own `WORK_DIR` and `PORT`:
`NETWORK=testnet LN_BACKEND_TYPE=LDK ESPLORA_SERVER=https://mempool.space/testnet/api`.

## Block explorer and chain data

With `NETWORK=testnet` and no `MEMPOOL_API` set, the hub uses `https://mempool.space/testnet/api`, so
every explorer link (transactions, addresses, channel funding transactions, on-chain receive, swaps)
opens `https://mempool.space/testnet/...`. `signet` uses `mempool.space/signet`, `testnet4` uses
`mempool.space/testnet4`, and mainnet is unchanged. An explicit `MEMPOOL_API` always wins (for
mutinynet: `NETWORK=signet MEMPOOL_API=https://mutinynet.com/api`). Every page shows a "You are on
testnet" notice, and Buy Bitcoin explains that it is mainnet-only.

## Funding

- Send testnet3 coins to the hub's on-chain address (Node → On-chain, or Receive → On-chain).
  Public faucets: https://bitcoinfaucet.uo1.net, https://coinfaucet.eu/en/btc-testnet/,
  https://testnet.help/en/btcfaucet/testnet. All of them use a captcha or bot check, so a person
  has to request the coins.
- Greenlight shows on-chain funds only after one confirmation; LDK shows them while unconfirmed.
- Opening a channel needs about 13,000 sat or more: CLN refuses channels whose usable capacity is
  under 10,000 sat after reserves and fees. 25,000–50,000 sat is comfortable.
- Receiving without a channel uses Blockstream's testnet LSP (LSPS2 JIT, 2,500 sat opening fee,
  minimum 2,501 sat). The payer needs a route to that LSP.

## Expected results and known limits (2026-11-21, UTC)

| Area | Greenlight | LDK |
|---|---|---|
| Setup, unlock, restart, backup export, restore into a fresh hub | works | works |
| Network label, explorer links, on-chain address, balances | works | works |
| On-chain send and receive | works | works |
| Invoices, lookup, list, NWC apps, budgets, pay-only scopes | works | works |
| Sign message | not supported by design (Greenlight signer, issue #739) | works |
| Hold invoices over NWC | not offered | works |
| Lightning payments | blocked on the hosted side, see below | works (keysend tested) |

Hosted-side blockers seen on testnet:
- **Old nodes start far behind the chain.** A node that was idle for weeks replays blocks at about
  5 blocks/s and reports "node syncing" until it reaches the tip (about 200,000 blocks, roughly 10
  hours, only while the hub is running). It cannot route payments meanwhile. A freshly registered
  node starts at the tip.
- **`Xpay` returns "no response from lightningd"** (Blockstream/greenlight#751). It still
  reproduces on a node that is syncing. The hub keeps such a payment pending and reconciles it later,
  so the budget is not released early.
- **The testnet LSP has no public channels** reachable from the faucet node's part of the graph, and
  it refuses channels below 5 BTC, so a JIT receive from another testnet wallet usually finds no
  route.

## Steps that still need funds

1. Fund a freshly registered Greenlight hub with at least 20,000 testnet sat.
2. Open a channel from it to a well-connected testnet node, for example the faucet node
   `039e7672537f2c742141dfadfd1b2e9a4871f9c064873347ab0fc0671f73f093f8@172.96.9.5:9735`.
3. Once the channel is active: pay an invoice, receive, send a keysend, pay through a budgeted NWC
   connection, and close the channel. A synced node is the real test of whether #751 still happens.
