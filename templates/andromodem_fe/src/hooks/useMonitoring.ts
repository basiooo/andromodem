import { useEffect } from "react"
import { toast } from "react-toastify"
import useSWR from "swr"
import useSWRMutation from "swr/mutation"

import { monitoringApi } from "@/api/monitoringApi"
import { useMonitoringConfigStore } from "@/stores/useMonitoringConfigStore"
import { type Device,DeviceState } from "@/types/device"
import type { DeviceMonitoringResponse } from "@/types/monitoring"

const useMonitoring = (device: Device | null) => {
  const { monitoringConfig, setMonitoringConfig, clearMonitoringConfig } = useMonitoringConfigStore()

  const {
    data: monitoringConfigData,
    mutate: mutateMonitoringConfig,
    isLoading: isLoadingMonitoringConfig,
    isValidating: isValidatingMonitoringConfig,
    error: monitoringConfigError
  } = useSWR<DeviceMonitoringResponse>(
    device?.state === DeviceState.ONLINE ? `monitoring_info_${device.serial}` : null,
    () => monitoringApi.getDeviceMonitoringConfiguration(device?.serial || ''),
    {
        revalidateIfStale: false,
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
        errorRetryCount: 0
    }
  )

  const {
    trigger: startMonitoring,
    isMutating: isStarting,
    error: startMonitoringError
  } = useSWRMutation(
    device?.state === DeviceState.ONLINE ? `start_monitoring_${device.serial}` : null,
    async () => {
      if (!device?.serial) throw new Error("Device serial not available")
      const result = await monitoringApi.startDeviceMonitoring(device.serial)
      await mutateMonitoringConfig()
      return result
    },
    {
      onSuccess: () => {
        toast.success("Monitoring started successfully!")
      },
      onError: (error) => {
        console.error('Failed to start monitoring:', error)
        const errorMessage = error?.response?.data?.message || error?.message || 'Failed to start monitoring'
        toast.error(errorMessage)
      }
    }
  )

  const {
    trigger: stopMonitoring,
    isMutating: isStopping,
    error: stopMonitoringError
  } = useSWRMutation(
    device?.state === DeviceState.ONLINE ? `stop_monitoring_${device.serial}` : null,
    async () => {
      if (!device?.serial) throw new Error("Device serial not available")
      const result = await monitoringApi.stopDeviceMonitoring(device.serial)
      await mutateMonitoringConfig()
      return result
    },
    {
      onSuccess: () => {
        toast.success("Monitoring stopped successfully!")
      },
      onError: (error) => {
        console.error('Failed to stop monitoring:', error)
        const errorMessage = error?.response?.data?.message || error?.message || 'Failed to stop monitoring'
        toast.error(errorMessage)
      }
    }
  )

  const {
    trigger: clearLogs,
    isMutating: isClearing,
    error: clearLogsError
  } = useSWRMutation(
    device?.state === DeviceState.ONLINE ? `clear_logs_${device.serial}` : null,
    async () => {
      if (!device?.serial) throw new Error("Device serial not available")
      return await monitoringApi.clearDeviceMonitoringLogs(device.serial)
    },
    {
      onSuccess: () => {
        toast.success("Logs cleared successfully!")
      },
      onError: (error) => {
        console.error('Failed to clear logs:', error)
        const errorMessage = error?.response?.data?.message || error?.message || 'Failed to clear logs'
        toast.error(errorMessage)
      }
    }
  )

  useEffect(() => {
    if (device === null) {
      clearMonitoringConfig()
    }
  }, [device, clearMonitoringConfig])

  useEffect(() => {
    if (monitoringConfigData && !isLoadingMonitoringConfig && monitoringConfigError === undefined && !isValidatingMonitoringConfig) {
      setMonitoringConfig(monitoringConfigData.data)
    }
    if (monitoringConfigError) {
      if (monitoringConfigError?.response?.status === 404) {
        clearMonitoringConfig()
      }else{
        toast.error(monitoringConfigError.message)
      }
    }
  }, [monitoringConfigData, isLoadingMonitoringConfig, monitoringConfigError, isValidatingMonitoringConfig, setMonitoringConfig, clearMonitoringConfig])

  return {
    monitoringConfig,
    monitoringConfigData,
    mutateMonitoringConfig,
    isLoadingMonitoringConfig,
    isValidatingMonitoringConfig,
    monitoringConfigError,
    
    startMonitoring,
    stopMonitoring,
    clearLogs,
    
    isStarting,
    isStopping,
    isClearing,
    
    startMonitoringError,
    stopMonitoringError,
    clearLogsError
  }
}

export default useMonitoring