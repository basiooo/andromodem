import {apiClient} from "@/api/apiClient"
import type { DeviceFeatureAvailabilityResponse, DeviceInfoResponse } from "@/types/device"

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