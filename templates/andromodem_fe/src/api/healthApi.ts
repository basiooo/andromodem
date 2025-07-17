import {apiClient} from "@/api/apiClient"
import type { BaseResponse } from "@/types/response"

export const HealthApi = {
    ping: async () => {
        const {data} = await apiClient.get<BaseResponse>(`/health/ping`)
        return data
    }
}