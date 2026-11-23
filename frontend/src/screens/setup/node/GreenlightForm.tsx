import { wordlist } from "@scure/bip39/wordlists/english.js";
import { useEffect } from "react";
import { useNavigate } from "react-router";

import * as bip39 from "@scure/bip39";
import Loading from "src/components/Loading";
import useSetupStore from "src/state/SetupStore";

/**
 * Greenlight needs no per-node configuration: the hosted node and local
 * signer are provisioned from the mnemonic at setup. Mirror LDKForm: set the
 * backend, generate a phrase if needed, and continue to security. The
 * backend-specific caveats are shown once in SetupSecurity.
 */
export function GreenlightForm() {
  const navigate = useNavigate();

  useEffect(() => {
    if (!useSetupStore.getState().nodeInfo.mnemonic) {
      useSetupStore.getState().updateNodeInfo({
        mnemonic: bip39.generateMnemonic(wordlist, 128),
      });
    }
    useSetupStore.getState().updateNodeInfo({
      backendType: "GREENLIGHT",
    });
    navigate("/setup/greenlight-cert", {
      replace: true,
    });
  }, [navigate]);

  return (
    <>
      <title>Loading... · Alby Hub</title>
      <Loading />
    </>
  );
}
