import { useEffect } from 'react'
import { toast } from 'react-toastify'
import useSWR from 'swr'

import { deviceApi } from '@/api/deviceApi'
import { useDevicesStore } from '@/stores/devicesStore'
import { type Device, type DeviceFeatureAvailabilityResponse, DeviceState } from '@/types/device'

const useFeatureAvailability = (device: Device | null) => {
  const { 
    setDeviceFeatureAvailabilities, 
    deviceFeatureAvailabilities, 
    unsetDeviceFeatureAvailabilities 
  } = useDevicesStore()

  const { data, isLoading, isValidating, error } = useSWR<DeviceFeatureAvailabilityResponse>(
    device?.state === DeviceState.ONLINE ? `device_feature_${device.serial}` : null,
    () => deviceApi.getDeviceFeatureAvailabilities(device?.serial || ''),
    {
      refreshInterval: 10000,
      revalidateIfStale: false,
      revalidateOnFocus: false,
      dedupingInterval: 55000,
      errorRetryInterval: 120000,
      errorRetryCount: 2
    }
  )

  useEffect(() => {
    if (device === null) {
      unsetDeviceFeatureAvailabilities()
    }
  }, [device, unsetDeviceFeatureAvailabilities])

  useEffect(() => {
    if (data && !isLoading && error === undefined && !isValidating) {
      setDeviceFeatureAvailabilities(data.data.feature_availabilities)
    }
    if (error) {
      toast.error(error.message)
    }
  }, [data, isLoading, error, isValidating, setDeviceFeatureAvailabilities, unsetDeviceFeatureAvailabilities])

  return {
    deviceFeatureAvailabilities,
    isLoading,
    isValidating,
    error
  }
}

export default useFeatureAvailability