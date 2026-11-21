import { FlaskConicalIcon } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "src/components/ui/alert";
import { useInfo } from "src/hooks/useInfo";
import { isMainnet } from "src/utils/isMainnet";

// Shown on every page when the hub is not running on bitcoin mainnet, so test
// coins are never mistaken for real ones.
export function TestNetworkAlert() {
  const { data: info } = useInfo();
  if (!info || isMainnet(info.network)) {
    return null;
  }
  return (
    <Alert variant="warning">
      <FlaskConicalIcon />
      <AlertTitle>You are on {info.network}</AlertTitle>
      <AlertDescription>
        Coins on this network have no value. Block explorer links open{" "}
        {info.mempoolUrl}.
      </AlertDescription>
    </Alert>
  );
}
