import { useCallback, useEffect, useRef, useState } from "react"
import { toast } from "react-toastify"

import { config } from "@/config"
import { useMonitoringStore } from "@/stores/monitoringStore"
import type { Device } from "@/types/device"
import type { MonitoringLog } from "@/types/monitoring"

const useListenMonitoringLog = ({ device }: { device: Device }) => {
  const eventSourceRef = useRef<EventSource | null>(null)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const retryCountRef = useRef<number>(0) 
  const [isConnected, setIsConnected] = useState<boolean>(false)
  
  const { logs, addLog, clearLogs: clearStoreLog } = useMonitoringStore()

  const MAX_RETRIES = 5
  const RETRY_INTERVAL_MS = 5000

  const clearLogs = useCallback(() => {
    clearStoreLog()
  }, [clearStoreLog])

  const connect = useCallback(() => {
    if (!device) {
      return
    }

    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }

    const source = new EventSource(`${config.BASE_URL}/event/devices/${device.serial}/monitoring/logs`)
    
    eventSourceRef.current = source

    source.addEventListener("open", () => {
      setIsConnected(true)
      console.log(`SSE connection established for monitoring logs: ${device.serial}`)
      retryCountRef.current = 0
    })

    source.addEventListener("message", (event: MessageEvent) => {
      try {
        if (event.data === "connected") {
          setIsConnected(true)
          return
        }
        
        const logData: MonitoringLog = JSON.parse(event.data)
        addLog(logData)
        setIsConnected(true)
      } catch (err) {
        console.error("Failed to parse monitoring log event:", err)
      }
    })

    source.addEventListener("error", (err: Event) => {
      console.warn(`SSE error for monitoring logs (${device.serial}):`, err)
      setIsConnected(false)

      source.close()

      if (retryCountRef.current < MAX_RETRIES) {
        retryCountRef.current++
        toast.warn(`Lost monitoring log connection, retrying... (${retryCountRef.current})`, {
          toastId: `monitoring_log_retry_${device.serial}_${retryCountRef.current}`
        })
        reconnectTimerRef.current = setTimeout(() => {
          console.log(`Retry ${retryCountRef.current} to connect to monitoring log SSE for ${device.serial}`)
          connect()
        }, RETRY_INTERVAL_MS)
      } else {
        toast.error(`Max retries reached for monitoring logs (${device.serial}). Please refresh the page.`, {
          toastId: `monitoring_log_max_retries_${device.serial}`
        })
      }
    })
  }, [device, addLog])

  const disconnect = () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
      eventSourceRef.current = null
    }
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current)
      reconnectTimerRef.current = null
    }
    setIsConnected(false)
    retryCountRef.current = 0
  }

  useEffect(() => {
    clearLogs()
    
    if (device) {
      connect()
    } else {
      disconnect()
    }

    return () => {
      disconnect()
    }
  }, [device, clearLogs, connect])

  return { logs, clearLogs, isConnected }
}

export default useListenMonitoringLog