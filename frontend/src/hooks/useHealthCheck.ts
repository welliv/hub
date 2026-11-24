import useSWR, { SWRConfiguration } from "swr";

import { isNodeSyncingAlarm } from "src/lib/health";
import { HealthResponse } from "src/types";
import { swrFetcher } from "src/utils/swr";

const pollConfiguration: SWRConfiguration<HealthResponse> = {
  refreshInterval: (latestData) =>
    latestData?.alarms?.some(isNodeSyncingAlarm) ? 10 * 1000 : 5 * 60 * 1000,
};

export function useHealthCheck(poll = true) {
  return useSWR<HealthResponse>(
    "/api/health",
    swrFetcher,
    poll ? pollConfiguration : undefined
  );
}
