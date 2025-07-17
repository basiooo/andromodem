import type { NetworkInfoResponse } from "@/types/network"
import type { BaseResponse } from "@/types/response"

import {apiClient} from "@/api/apiClient"

export const networkApi = {
    getNetworkInfo: async (serial: string) => {
        const {data} = await apiClient.get<NetworkInfoResponse>(`/devices/${serial}/network`)
        return data
    },
    toggleAirplaneMode: async (serial: string) => {
        const {data} = await apiClient.post<BaseResponse>(`/devices/${serial}/network/airplane-mode`)
        return data
    },
    toggleMobileData: async (serial: string) => {
        const {data} = await apiClient.post<BaseResponse>(`/devices/${serial}/network/mobile-data`)
        return data
    }
}