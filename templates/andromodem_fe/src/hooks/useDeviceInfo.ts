import { useEffect, useState } from 'react'
import { toast } from 'react-toastify'
import useSWR from 'swr'

import { deviceApi } from '@/api/deviceApi'
import { type Device, type DeviceInfo, type DeviceInfoResponse, DeviceState } from '@/types/device'

const useDeviceInfo = (device: Device | null) => {
  const [deviceInfo, setDeviceInfo] = useState<DeviceInfo | null>(null)
  const [upTimeSecond, setUpTimeSecond] = useState<number>(0)

  const { data, mutate, isLoading, isValidating, error } = useSWR<DeviceInfoResponse>(
    device?.state === DeviceState.ONLINE ? `device_info_${device.serial}` : null,
    () => deviceApi.getDeviceInfo(device?.serial || ''),
    {
      revalidateIfStale: false,
      revalidateOnFocus: false
    }
  )

  useEffect(() => {
    if (device === null) {
      setDeviceInfo(null)
    }
  }, [device])

  useEffect(() => {
    if (upTimeSecond > 0) {
      const interval = setInterval(() => {
        setUpTimeSecond((prev) => prev + 1)
      }, 1000)
      return () => clearInterval(interval)
    }
  }, [upTimeSecond])

  useEffect(() => {
    if (data && !isLoading && error === undefined && !isValidating) {
      setDeviceInfo(data.data)
      setUpTimeSecond(data.data.other.uptime_second || 0)
    }
    if (error) {
      toast.error(error.message)
    }
  }, [data, isLoading, error, isValidating])

  return {
    deviceInfo,
    upTimeSecond,
    isLoading,
    isValidating,
    error,
    mutate
  }
}

export default useDeviceInfo