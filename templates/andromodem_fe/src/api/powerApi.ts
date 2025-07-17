import {apiClient} from "@/api/apiClient"
import type { DevicePowerActionType } from "@/types/power"
import type { BaseResponse } from "@/types/response"

export const PowerApi = {
    powerAction: async (serial: string, action: DevicePowerActionType) => {
        const {data} = await apiClient.post<BaseResponse>(`/devices/${serial}/power`, {action})
        return data
    }
}