import { ReactElement } from "react";
import { GreenlightIcon } from "src/components/icons/Greenlight";
import { LDKIcon } from "src/components/icons/LDK";
import { PhoenixdIcon } from "src/components/icons/Phoenixd";
import { BackendType } from "src/types";

import barkDark from "src/assets/images/node/bark-dark.svg";
import barkLight from "src/assets/images/node/bark-light.svg";
import cashu from "src/assets/images/node/cashu.png";
import cln from "src/assets/images/node/cln.png";
import lnd from "src/assets/images/node/lnd.png";

type BackendTypeConfig = {
  title: string;
  subtitle?: string;
  icon: ReactElement;
  hasMnemonic: boolean;
  hasChannelManagement: boolean;
  hasNodeBackup: boolean;
};

export const backendTypeConfigs: Record<BackendType, BackendTypeConfig> = {
  LDK: {
    title: "LDK",
    icon: <LDKIcon />,
    hasMnemonic: true,
    hasChannelManagement: true,
    hasNodeBackup: true,
  },
  PHOENIX: {
    title: "phoenixd",
    icon: <PhoenixdIcon />,
    hasMnemonic: false,
    hasChannelManagement: false,
    hasNodeBackup: false,
  },
  LND: {
    title: "LND",
    icon: <img src={lnd} />,
    hasMnemonic: false,
    hasChannelManagement: true,
    hasNodeBackup: false,
  },
  CASHU: {
    title: "Cashu Mint",
    icon: <img src={cashu} />,
    hasMnemonic: true,
    hasChannelManagement: false,
    hasNodeBackup: false,
  },
  CLN: {
    title: "CLN",
    icon: <img src={cln} />,
    hasMnemonic: false,
    hasChannelManagement: true,
    hasNodeBackup: false,
  },
  GREENLIGHT: {
    title: "Greenlight",
    subtitle: "Experimental",
    // Official two-tone emblem on a dark badge. The emblem's white half needs
    // a dark ground; the yellow half is Blockstream's primary Greenlight yellow.
    icon: <GreenlightIcon />,
    // 12-word BIP-39 is both GL node seed and Settings backup phrase
    hasMnemonic: true,
    // Real CLN channels on Blockstream cloud — not Bark stubs
    hasChannelManagement: true,
    // The Hub .bkp includes device credentials and the signer's static
    // channel backup, so Greenlight also has a real node backup.
    hasNodeBackup: true,
  },
  BARK: {
    title: "Bark",
    icon: (
      <>
        <img src={barkLight} className="dark:hidden" />
        <img src={barkDark} className="hidden dark:block" />
      </>
    ),
    hasMnemonic: true,
    hasChannelManagement: false,
    hasNodeBackup: false,
  },
};
