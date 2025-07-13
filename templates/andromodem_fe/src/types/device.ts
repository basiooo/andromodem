import type { BaseResponse } from "./response"


export const DeviceState = {
  UNKNOWN: "Unknown",
  ONLINE: "Online",
  DISCONNECT: "Disconnected",
  OFFLINE: "Offline",
  UNAUTHORIZED: "Unauthorized",
  AUTHORIZING: "Authorizing",
  Recovery: "Recovery"
} as const

export type DeviceStateValue = typeof DeviceState[keyof typeof DeviceState];


export interface Device {
  serial: string;
  model: string;
  product: string;
  state: DeviceStateValue;
  old_state: DeviceStateValue;
  new_state: DeviceStateValue;
  android_version: string;
}

export type Devices = Device[];

export interface DeviceStorage {
  data_total: number;
  data_free: number;
  data_used: number;
  system_total: number;
  system_free: number;
  system_used: number;
};

export interface DeviceMemory {
  mem_total: number;
  mem_free: number;
  mem_used: number;
  swap_total: number;
  swap_free: number;
  swap_used: number;
};

export interface DeviceBattery {
  ac_powered: boolean;
  usb_powered: boolean;
  wireless_powered: boolean;
  max_charging_current: number;
  max_charging_voltage: number;
  charge_counter: number;
  status: string;
  health: string;
  present: boolean;
  level: number;
  scale: number;
  temperature: number;
  technology: string;
}


export interface DeviceInfo {
  prop: {
    model: string;
    brand: string;
    name: string;
    android_version: string;
    fingerprint: string;
    sdk: string;
    security_patch: string;
    processor: string;
    abi: string;
  };
  root: {
    is_rooted: boolean;
    super_user_allow_shell_access: boolean;
    version: string;
    name: string;
  };
  battery: DeviceBattery;
  other: {
    uptime: string;
    uptime_second: number;
    busybox_installed: boolean;
    kernel_version: string;
  };
  memory: DeviceMemory;
  storage: DeviceStorage
};

export interface DeviceInfoResponse extends BaseResponse{
  data: DeviceInfo
}

export interface FeatureAvailability{
  feature: string;
  key: string;
  available: boolean;
  message: string;
}
export type FeatureAvailabilities = FeatureAvailability[];


export interface DeviceFeatureAvailabilityResponse extends BaseResponse{
  data: {
    feature_availabilities: FeatureAvailabilities;
  }
}