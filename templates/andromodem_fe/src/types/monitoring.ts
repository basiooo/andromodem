import type { BaseResponse } from "./response"

export const MonitoringMethod = {
  HTTP :'http',
  HTTPS : 'https',
  WS : 'ws',
  PING : 'icmp',
  PING_BY_DEVICE : 'ping_by_device'
} as const

export type MonitoringMethodValue = (typeof MonitoringMethod)[keyof typeof MonitoringMethod];

export interface MonitoringConfig {
  serial: string;
  host: string;
  method: MonitoringMethodValue;
  max_failures: number;
  checking_interval: number;
  airplane_mode_delay: number;
  is_active: boolean;
  created_at: Date; 
  updated_at: Date;
}

export interface DeviceMonitoringResponse extends BaseResponse {
    data: MonitoringConfig;
}

export interface MonitoringConfigPayload {
  host: string;
  method: MonitoringMethodValue;
  max_failures: number;
  airplane_mode_delay: number;
  checking_interval: number;
}

export interface MonitoringLog{
  serial: string,
  success:boolean,
  message:string,
  timestamp:Date,
}