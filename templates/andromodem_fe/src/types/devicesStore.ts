import type { Device, Devices,FeatureAvailabilities } from "@/types/device"

export interface DevicesStore {
    devices: Devices;
    deviceUsed: Device | null;
    deviceFeatureAvailabilities: FeatureAvailabilities | null;
    setDevices(devices: Devices): void;
    addDevice(device: Device): void;
    addOrUpdateDevices(serial: string, newDevice: Device): void;
    setDeviceUsed(deviceUsed: Device): void;
    unsetDeviceUsed(): void;
    setDeviceFeatureAvailabilities(featureAvailabilities: FeatureAvailabilities): void;
    unsetDeviceFeatureAvailabilities(): void;
}