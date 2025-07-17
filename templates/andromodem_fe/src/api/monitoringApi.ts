import { apiClient } from "@/api/apiClient"
import type { DeviceMonitoringResponse, MonitoringConfigPayload } from "@/types/monitoring"
import type { BaseResponse } from "@/types/response"

export const monitoringApi = {
    getDeviceMonitoringConfiguration: async (serial: string)=>{
        const {data} = await apiClient.get<DeviceMonitoringResponse>(`/devices/${serial}/monitoring`)
        return data
    },
    createDeviceMonitoringConfiguration: async (serial: string, payload: MonitoringConfigPayload )=>{
        const {data} = await apiClient.post<DeviceMonitoringResponse>(`/devices/${serial}/monitoring`,payload)
        return data
    },
    updateDeviceMonitoringConfiguration: async (serial: string, payload: MonitoringConfigPayload )=>{
        const {data} = await apiClient.put<DeviceMonitoringResponse>(`/devices/${serial}/monitoring`,payload)
        return data
    },
    stopDeviceMonitoring: async (serial: string)=>{
        const {data} = await apiClient.post<BaseResponse>(`/devices/${serial}/monitoring/stop`)
        return data
    },
    startDeviceMonitoring: async (serial: string)=>{
        const {data} = await apiClient.post<BaseResponse>(`/devices/${serial}/monitoring/start`)
        return data
    },
    clearDeviceMonitoringLogs: async (serial: string)=>{
        const {data} = await apiClient.delete<BaseResponse>(`/devices/${serial}/monitoring/logs`)
        return data
    }
}