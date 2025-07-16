import type { BaseResponse } from "./response"

export type DevicePowerActionType = 
  | "reboot"
  | "reboot_recovery"
  | "reboot_bootloader"
  | "power_off";

export interface PowerActionRequest extends BaseResponse {
  action: DevicePowerActionType;
}