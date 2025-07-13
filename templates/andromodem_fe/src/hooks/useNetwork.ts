import { useEffect, useState } from 'react'
import { toast } from 'react-toastify'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'

import { networkApi } from '@/api/networkApi'
import { type Device,DeviceState } from '@/types/device'
import type { NetworkData, NetworkInfoResponse } from '@/types/network'

const useNetwork = (device: Device | null) => {
  const [networkInfo, setNetworkInfo] = useState<NetworkData | null>(null)

  const {
    data: networkData,
    mutate: mutateNetworkInfo,
    isLoading: isLoadingNetworkInfo,
    isValidating: isValidatingNetworkInfo,
    error: networkInfoError
  } = useSWR<NetworkInfoResponse>(
    device?.state === DeviceState.ONLINE ? `network_info_${device.serial}` : null,
    () => networkApi.getNetworkInfo(device?.serial || ''),
    {
      revalidateIfStale: false,
      revalidateOnFocus: false
    }
  )

  const {
    trigger: toggleMobileData,
    isMutating: isTogglingMobileData,
    error: toggleMobileDataError
  } = useSWRMutation(
    device?.state === DeviceState.ONLINE ? `toggle_mobile_data_${device.serial}` : null,
    async () => {
      if (!device?.serial) throw new Error('Device serial is required')
      const result = await networkApi.toggleMobileData(device.serial)

      await mutateNetworkInfo()
      return result
    },
    {
      onSuccess: (data) => {
        toast.success(data.message)
      },
      onError: (error) => {
        if(error?.response?.data?.message){
          mutateNetworkInfo()
          toast.error(error.response.data.message)
        }else{
          toast.error(error.message)
        }
      }
    }
  )

  const {
    trigger: toggleAirplaneMode,
    isMutating: isTogglingAirplaneMode,
    error: toggleAirplaneModeError
  } = useSWRMutation(
    device?.state === DeviceState.ONLINE ? `toggle_airplane_mode_${device.serial}` : null,
    async () => {
      if (!device?.serial) throw new Error('Device serial is required')
      const result = await networkApi.toggleAirplaneMode(device.serial)

      await mutateNetworkInfo()
      return result
    },
    {
      onSuccess: (data) => {
        toast.success(data.message)
      },
      onError: (error) => {
        if(error?.response?.data?.message){
          mutateNetworkInfo()
          toast.error(error.response.data.message)
        }else{
          toast.error(error.message)
        }
      }
    }
  )


  useEffect(() => {
    if (device === null) {
      setNetworkInfo(null)
    }
  }, [device])

  useEffect(() => {
    if (networkData && !isLoadingNetworkInfo && networkInfoError === undefined && !isValidatingNetworkInfo) {
      setNetworkInfo(networkData.data)
    }
    if (networkInfoError) {
      toast.error(networkInfoError.message)
    }
  }, [networkData, isLoadingNetworkInfo, networkInfoError, isValidatingNetworkInfo])

  return {
    networkInfo,
    isLoadingNetworkInfo,
    isValidatingNetworkInfo,
    networkInfoError,
    mutateNetworkInfo,

    toggleMobileData,
    isTogglingMobileData,
    toggleMobileDataError,

    toggleAirplaneMode,
    isTogglingAirplaneMode,
    toggleAirplaneModeError
  }
}

export default useNetwork