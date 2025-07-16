import type { MonitoringConfig, MonitoringLog } from "@/types/monitoring"

export interface MonitoringStore {
    logs: MonitoringLog[];
    addLog(log: MonitoringLog): void;
    clearLogs(): void;
}

export interface MonitoringConfigStore {
  monitoringConfig: MonitoringConfig | null;
  setMonitoringConfig: (config: MonitoringConfig | null) => void;
  clearMonitoringConfig: () => void;
}