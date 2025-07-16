import { create } from "zustand"

import type { Device, Devices,FeatureAvailabilities } from "@/types/device"
import type { DevicesStore } from "@/types/devicesStore"


export const useDevicesStore = create<DevicesStore>((set) => {
    return {
        devices: [],
        deviceUsed: null,
        deviceFeatureAvailabilities: null,
        setDevices: (devices: Devices) => {
            set({
                devices
            })
        },
        addDevice: (device: Device) => {
            set((state) => ({
                devices: [...state.devices, device]
            }))
        },
        addOrUpdateDevices: (serial: string, newDevice: Device) => set((state) => {
            const deviceIndex = state.devices.findIndex((device: Device) => device.serial === serial)
            let result = null
            if (deviceIndex !== -1) {
                const updatedDevices = [...state.devices]
                updatedDevices[deviceIndex] = { ...updatedDevices[deviceIndex], ...newDevice }
                result = { devices: updatedDevices }
            } else {
                result = { devices: [...state.devices, newDevice] }
            }
            if (state.deviceUsed && state.deviceUsed.serial === serial) {
                state.setDeviceUsed(newDevice)
            }
            return result
        }),
        setDeviceUsed: (deviceUsed: Device) => {
            set({
                deviceUsed: deviceUsed
            })
        },
        unsetDeviceUsed: () => {
            set({
                deviceUsed: null
            })
        },
        setDeviceFeatureAvailabilities: (featureAvailabilities: FeatureAvailabilities) => {
            set({
                deviceFeatureAvailabilities: featureAvailabilities
            })
        },
        unsetDeviceFeatureAvailabilities: () => {
            set({
                deviceFeatureAvailabilities: [] 
            })
        }
    }
})