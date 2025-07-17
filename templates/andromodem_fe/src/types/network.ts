import type { BaseResponse } from "@/types/response"

export interface NetworkData {
  airplane_mode: boolean;
  ip_routes: IpRoute[];
  apn: Apn;
  sims: SimInfo[];
}

export interface IpRoute {
  interface: string;
  ip: string;
}

export interface Apn {
  name: string;
  apn: string;
  proxy: string;
  port: string;
  mms_proxy: string;
  mms_port: string;
  username: string;
  password: string;
  server: string;
  mcc: string;
  mnc: string;
  type: string;
  protocol: string;
}


export interface SimInfo {
  name: string;
  connection_state: MobileDataConnectionStateValue;
  sim_slot: number;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  signal_strength: any;
}

export interface NetworkInfoResponse extends BaseResponse{
    data: NetworkData
}

export const MobileDataConnectionState = {
    UNKNOWN: "Unknown",
    DISCONNECTED: "Disconnected",
    CONNECTING: "Connecting",
    CONNECTED: "Connected",
    SUSPENDED: "Suspended"
} as const

export type MobileDataConnectionStateValue = typeof MobileDataConnectionState[keyof typeof MobileDataConnectionState];
