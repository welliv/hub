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

## Results of the testnet run (2026-11-21/22, UTC)

| Area | Greenlight | LDK |
|---|---|---|
| Setup, unlock, restart, backup export, restore into a fresh hub | works | works |
| Network label, explorer links, on-chain address, balances | works | works |
| On-chain send and receive, send-all sweep | works (CLN keeps a 546 sat emergency reserve while an anchor channel exists) | works |
| Invoices, lookup, list, NWC apps, budgets, pay-only scopes (RESTRICTED) | works | works |
| NWC over-budget payment → QUOTA_EXCEEDED | works | works |
| NWC notifications | subscribed, no payment to notify | works (payment_sent delivered) |
| Sign message | not supported by design (Greenlight signer, issue #739) | works |
| Hold invoices over NWC | not offered | works |
| Open channel (100k to a testnet CLN node) | funding confirms, both sides send channel_ready, but the hosted node stays in CHANNELD_AWAITING_LOCKIN (see below) | works (30k, active after 3 confirmations) |
| Cooperative close | works (from AWAITING_LOCKIN) | waits on an in-flight HTLC |
| Pay invoice / keysend / NWC pay | not testable: no locked-in channel; on a syncing node Xpay fails with "no response from lightningd" (#751) | keysend works; invoice payments work when a route exists |
| JIT receive (LSPS2, LSP 034066…) | invoice created, 2,500 sat fee | invoice created (testnet LSPS2 default now built in), 2,500 sat fee, min 2,501 sat |
| JIT payment | no route: the LSP has no public channels on testnet | same |

Hosted-side blockers seen on testnet:
- **Channel never locks in.** A 100k channel from a freshly registered node confirmed (13+
  confirmations, the wallet sees it), the peer sent `channel_ready` and the node's own channeld
  reports "Channel ready for use", yet lightningd stayed in `CHANNELD_AWAITING_LOCKIN` and kept
  re-broadcasting the funding transaction ("Transaction outputs already in utxo set"). Restarting
  the hosted node twice did not help. Before the first restart the peer never answered the
  channel reestablish. The hub now shows such a channel as opening instead of logging
  "Could not find status indicating number of confirmations required".
- **Old nodes start far behind the chain.** A node that was idle for weeks replays blocks at about
  5 blocks/s and reports "node syncing" until it reaches the tip (about 200,000 blocks, roughly 10
  hours, only while the hub is running). It cannot route payments meanwhile.
- **`Xpay` returns "no response from lightningd"** (Blockstream/greenlight#751) on a syncing node.
  Not yet tested on a synced node with an active channel.
- **No inbound connections.** A Greenlight node is not announced, so another node cannot open a
  channel to it; channels to it come from its own opens or from the LSP.
- **The testnet LSP has no public channels**, and refuses channels below 5 BTC, so JIT receives
  find no route. ACINQ's testnet node requires 400,000 sat per channel.

## Steps that still need funds

1. Fund a freshly registered Greenlight hub with at least 20,000 testnet sat (more for peers with
   high minimums).
2. Open a channel from it to a well-connected testnet node, for example the faucet node
   `039e7672537f2c742141dfadfd1b2e9a4871f9c064873347ab0fc0671f73f093f8@172.96.9.5:9735`.
3. Once the channel is active: pay an invoice, receive, send a keysend, pay through a budgeted NWC
   connection, and close the channel. A synced node is the real test of whether #751 still happens.
