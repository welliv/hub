import { HealthAlarm } from "src/types";

export type NodeHealthAlarmDetails = {
  isReady?: boolean;
  syncing?: boolean;
  not_ready_reason?: string;
  node_connected?: boolean;
  last_error?: string;
  last_check_at?: number;
  signer?: {
    running?: boolean;
    last_error?: string;
  };
  internalNodeStatus?: {
    ready?: boolean;
    syncing?: boolean;
    not_ready_reason?: string;
    node_connected?: boolean;
    last_error?: string;
    last_check_at?: number;
    signer?: {
      running?: boolean;
      last_error?: string;
    };
  };
};

export function getNodeNotReadyDetails(
  alarm: HealthAlarm
): NodeHealthAlarmDetails | undefined {
  if (alarm.kind !== "node_not_ready") {
    return undefined;
  }
  return alarm.rawDetails as NodeHealthAlarmDetails | undefined;
}

export function getNodeNotReadyReason(alarm: HealthAlarm): string | undefined {
  const details = getNodeNotReadyDetails(alarm);
  return (
    details?.internalNodeStatus?.not_ready_reason ?? details?.not_ready_reason
  );
}

export function isNodeSyncingAlarm(alarm: HealthAlarm): boolean {
  const details = getNodeNotReadyDetails(alarm);
  return (
    details?.internalNodeStatus?.syncing === true ||
    details?.syncing === true ||
    getNodeNotReadyReason(alarm) === "node syncing"
  );
}
