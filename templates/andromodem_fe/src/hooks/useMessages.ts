import { useEffect, useMemo, useState } from 'react'
import { toast } from 'react-toastify'
import useSWR from 'swr'

import { MessageApi } from '@/api/messageApi'
import { type Device, DeviceState } from '@/types/device'
import type { Messages, MessagesResponse } from '@/types/message'

const useMessages = (device: Device | null) => {
  const [messages, setMessages] = useState<Messages>([])
  const [showValue, setShowValue] = useState<number>(10)

  const { data, mutate, isLoading, isValidating, error } = useSWR<MessagesResponse>(
    device?.state === DeviceState.ONLINE ? `message_inbox_${device?.serial}` : null,
    () => MessageApi.getMessages(device?.serial || ''),
    {
      revalidateIfStale: false,
      revalidateOnFocus: false,
      errorRetryCount: 0,
      revalidateOnReconnect: false
    }
  )

  const showOptions: number[] = [10, 20, 50, 100, -1]

  const finalMessages = useMemo(() => {
    const result = messages
    if (showValue === -1) {
      return result
    }
    return result.slice(0, showValue)
  }, [messages, showValue])

  useEffect(() => {
    if (data && !isLoading && error === undefined && !isValidating) {
      setMessages(data.data.messages ?? [])
    }
    if (error) {
      const error_data = error.response.data
      toast.error(`Messages: ${error_data.message}`, {
        toastId: 'message_list'
      })
      setMessages([])
    }
  }, [data, error, device, isLoading, isValidating])

  const handleShowOptionChange = (value: string) => {
    const val = parseInt(value)
    setShowValue(val)
  }

  return {
    messages,
    finalMessages,
    showValue,
    showOptions,
    isLoading,
    isValidating,
    error,
    mutate,
    handleShowOptionChange
  }
}

export default useMessages