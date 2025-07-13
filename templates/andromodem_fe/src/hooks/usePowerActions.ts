import { useState } from 'react'
import { toast } from 'react-toastify'
import useSWRMutation from 'swr/mutation'

import { PowerApi } from '@/api/powerApi'
import { type Device, DeviceState } from '@/types/device'
import type { DevicePowerActionType } from '@/types/power'
import type { BaseResponse } from '@/types/response'

const usePowerActions = (device: Device | null) => {
  const [action, setAction] = useState<DevicePowerActionType>('reboot')

  const powerDeviceHandler = async ([serial]: [string], { arg }: { arg: DevicePowerActionType }) => {
    return await PowerApi.powerAction(serial, arg)
  }

  const { trigger, isMutating: isLoading } = useSWRMutation<
    BaseResponse,
    Error,
    [string],
    DevicePowerActionType
  >([device?.serial || ''], powerDeviceHandler)

  const isDisabled = isLoading || device?.state !== DeviceState.ONLINE

  const powerActionHandler = (value: DevicePowerActionType) => {
    setAction(value)
    const power_modal = document.getElementById('power_modal') as HTMLDialogElement
    power_modal.showModal()
  }

  const powerActionHandlerConfirm = async () => {
    const close_power_modal = document.getElementById('close_power_modal') as HTMLButtonElement
    close_power_modal.click()
    try {
      const response = await trigger(action)
      toast.success(response.message, {
        toastId: `${device?.serial}_action`
      })
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message, {
          toastId: `${device?.serial}_action`
        })
      } else {
        toast.error('Unknown error occurred', {
          toastId: `${device?.serial}_action`
        })
      }
    }
  }

  return {
    action,
    isLoading,
    isDisabled,
    powerActionHandler,
    powerActionHandlerConfirm
  }
}

export default usePowerActions