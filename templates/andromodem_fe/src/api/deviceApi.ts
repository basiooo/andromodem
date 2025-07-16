import type { DeviceFeatureAvailabilityResponse, DeviceInfoResponse } from "@/types/device"

import {apiClient} from "./apiClient"

export const deviceApi = {
    getDeviceInfo: async (serial: string) => {
        const {data} = await apiClient.get<DeviceInfoResponse>(`/devices/${serial}`)
        return data
    },
    getDeviceFeatureAvailabilities: async (serial: string) => {
        const {data} = await apiClient.get<DeviceFeatureAvailabilityResponse>(`/devices/${serial}/feature-availabilities`)
        return data
    }
}