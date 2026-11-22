import useSWR from "swr";

import { RecommendedChannelPeer } from "src/types";
import { swrFetcher } from "src/utils/swr";

// Alby's public suggestion catalog has no testnet entries, so the on-chain
// peer picker is empty on testnet. These are well-connected testnet3 nodes
// (mempool.space connectivity ranking, checked 2026-11-22). They are onchain
// peers only: the LSPS1 lightning flow still has no testnet seller, so this
// does not invent one — it gives the user a working path to open a channel.
const PEER_ICON =
  "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E%3Ccircle cx='16' cy='16' r='14' fill='%2364748b'/%3E%3C/svg%3E";

const TESTNET_ONCHAIN_PEERS: RecommendedChannelPeer[] = [
  {
    network: "testnet",
    paymentMethod: "onchain",
    name: "1ML.com node ALPHA (testnet)",
    image: PEER_ICON,
    minimumChannelSizeSat: 100000,
    maximumChannelSizeSat: 16000000,
    note: "Well-connected testnet peer (~830 channels).",
    description: "",
    publicChannelsAllowed: true,
    pubkey:
      "02312627fdf07fbdd7e5ddb136611bdde9b00d26821d14d94891395452f67af248",
    host: "66.109.24.42:9735",
  },
  {
    network: "testnet",
    paymentMethod: "onchain",
    name: "aranguren.org (testnet)",
    image: PEER_ICON,
    minimumChannelSizeSat: 100000,
    maximumChannelSizeSat: 16000000,
    note: "High-capacity testnet routing node (~750 channels).",
    description: "",
    publicChannelsAllowed: true,
    pubkey:
      "038863cf8ab91046230f561cd5b386cbff8309fa02e3f0c3ed161a3aeb64a643b9",
    host: "203.132.94.196:9735",
  },
  {
    network: "testnet",
    paymentMethod: "onchain",
    name: "voltage-test (testnet)",
    image: PEER_ICON,
    minimumChannelSizeSat: 100000,
    maximumChannelSizeSat: 16000000,
    note: "Testnet node with ~450 channels.",
    description: "",
    publicChannelsAllowed: true,
    pubkey:
      "02cf71da3f277c2a30a348dfced77a9b7d81fb578c2b9117967100352051626b84",
    host: "54.214.32.132:20402",
  },
  {
    network: "testnet",
    paymentMethod: "onchain",
    name: "SATLOOT faucet (testnet)",
    image: PEER_ICON,
    minimumChannelSizeSat: 50000,
    maximumChannelSizeSat: 16000000,
    note: "Faucet peer used in the Greenlight testnet trial docs.",
    description: "",
    publicChannelsAllowed: true,
    pubkey:
      "039e7672537f2c742141dfadfd1b2e9a4871f9c064873347ab0fc0671f73f093f8",
    host: "172.96.9.5:9735",
  },
];

export function useChannelPeerSuggestions() {
  return useSWR<RecommendedChannelPeer[]>(
    "/api/channels/suggestions",
    async (key: string) => {
      const remote = (await swrFetcher(key)) as RecommendedChannelPeer[] | null;
      const knownPubkeys = new Set(
        (remote ?? []).map((peer) => peer.pubkey).filter(Boolean)
      );
      return [
        ...(remote ?? []),
        ...TESTNET_ONCHAIN_PEERS.filter(
          (peer) => !knownPubkeys.has(peer.pubkey)
        ),
      ];
    }
  );
}
