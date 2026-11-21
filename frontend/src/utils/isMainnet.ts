export function isMainnet(network: string | undefined) {
  return !network || network === "bitcoin" || network === "mainnet";
}
