import { useCallback, useEffect, useRef, useState } from "react"
import { toast } from "react-toastify"

import { config } from "@/config"
import { useDevicesStore } from "@/stores/devicesStore"
import type { Device,Devices } from "@/types/device"

const useListenDevices = (): { devices: Devices; isConnected: boolean } => {
  const { devices, addOrUpdateDevices } = useDevicesStore()
  const eventSourceRef = useRef<EventSource | null>(null)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const retryCountRef = useRef<number>(0) 
  const [isConnected, setIsConnected] = useState<boolean>(false)

  const MAX_RETRIES = 5
  const RETRY_INTERVAL_MS = 5000

  const connect = useCallback(() => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }

    const source = new EventSource(`${config.BASE_URL}/event/devices`)
    
    eventSourceRef.current = source

    source.addEventListener("open", () => {
      setIsConnected(true)
      console.log("SSE connection established")
      retryCountRef.current = 0
    })

    source.addEventListener("message", (event: MessageEvent) => {
      try {
        if (event.data === "connected") {
          setIsConnected(true)
          return
        }
        const data: Device = JSON.parse(event.data)
        toast.info(`Device "${data.serial}" ${data.new_state}`, {
          toastId: `${data.serial}_${data.new_state}`
        })
        addOrUpdateDevices(data.serial, data)
        setIsConnected(true)
      } catch (err) {
        console.error("Failed to parse device event:", err)
      }
    })

    source.addEventListener("error", (err: Event) => {
      console.warn("SSE error:", err)
      setIsConnected(false)

      source.close()

      if (retryCountRef.current < MAX_RETRIES) {
        retryCountRef.current++
        toast.warn(`Lost connection, retrying... (${retryCountRef.current})`, {
          toastId: `sse_retry_${retryCountRef.current}`
        })
        reconnectTimerRef.current = setTimeout(() => {
          console.log(`Retry ${retryCountRef.current} to connect to SSE`)
          connect()
        }, RETRY_INTERVAL_MS)
      } else {
        toast.error("Max retries reached. Please refresh the page.", {
          toastId: "sse_max_retries"
        })
      }
    })
  }, [addOrUpdateDevices])

  useEffect(() => {
    connect()

    return () => {
      eventSourceRef.current?.close()
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current)
      }
    }
  }, [connect])

  return { devices, isConnected }
}

export default useListenDevices
