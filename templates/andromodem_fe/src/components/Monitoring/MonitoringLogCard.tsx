import type { FC } from "react"
import { useEffect, useState } from "react"
import { toast } from "react-toastify"

import { monitoringApi } from "@/api/monitoringApi"
import useListenMonitoringLog from "@/hooks/useListenMonitoringLog"
import useMonitoring from "@/hooks/useMonitoring"
import type { Device } from "@/types/device"

const MonitoringLogCard: FC<{ device: Device }> = ({ device }) => {
  const { logs, clearLogs } = useListenMonitoringLog({ device })
  const [isClearing, setIsClearing] = useState(false)
  const { mutateMonitoringConfig, monitoringConfig } = useMonitoring(device)
  
  useEffect(() => {
    if (logs.length > 0 && monitoringConfig?.is_active && logs[logs.length - 1].message === "Monitoring task stopped") {
        mutateMonitoringConfig()
    }
  }, [logs, monitoringConfig?.is_active, mutateMonitoringConfig])

  const formatTimestamp = (timestamp: Date) => {
    const date = new Date(timestamp)
    const day = date.getDate().toString().padStart(2, '0')
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const year = date.getFullYear()
    const hours = date.getHours().toString().padStart(2, '0')
    const minutes = date.getMinutes().toString().padStart(2, '0')
    const seconds = date.getSeconds().toString().padStart(2, '0')
    
    return `${day}-${month}-${year}, ${hours}:${minutes}:${seconds}`
  }

  const handleClearLogs = async () => {
    try {
      setIsClearing(true)
      await monitoringApi.clearDeviceMonitoringLogs(device.serial)
      clearLogs()
      toast.success("Monitoring logs cleared successfully!")
    } 
    // eslint-disable-next-line
    catch (error: any) {
      console.error('Failed to clear monitoring logs:', error)
      const errorMessage = error?.response?.data?.message || error?.message || 'Failed to clear monitoring logs'
      toast.error(errorMessage)
    } finally {
      setIsClearing(false)
    }
  }

  return (
    <div className="card w-full bg-base-100 shadow-sm">
      <div className="card-body">
        <div className="flex justify-between items-center mb-4">
          <h2 className="card-title text-sm md:text-base flex items-center gap-2">
            Monitoring Logs
          </h2>
          <div className="flex items-center gap-4">
            <span className="text-xs md:text-sm text-base-content/70">
              {logs.length} entries
            </span>
            <button
              onClick={handleClearLogs}
              disabled={logs.length === 0 || isClearing}
              className="btn btn-warning btn-xs md:btn-sm"
            >
              {isClearing ? (
                <>
                  <span className="loading loading-spinner loading-xs"></span>
                  Clearing...
                </>
              ) : (
                <>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  Clear Logs
                </>
              )}
            </button>
          </div>
        </div>

        <div className="bg-base-100 rounded-lg border border-base-300 overflow-hidden">
          <div className="bg-base-300 px-4 py-2 flex items-center gap-2">
            <div className="flex gap-1">
              <div className="w-3 h-3 rounded-full bg-red-500"></div>
              <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
              <div className="w-3 h-3 rounded-full bg-green-500"></div>
            </div>
            <span className="text-base-content/70 text-sm ml-2">
              monitoring@{device.serial}
            </span>
          </div>

          <div 
            className="p-4 h-96 overflow-y-auto font-mono text-sm text-green-400 bg-base-200"
          >
            {logs.length === 0 ? (
              <div className="text-base-content/70">
                No logs yet
              </div>
            ) : (
              logs.map((log, index) => (
                <div key={`${log.timestamp}-${index}`} className="mb-2 pb-2 border-b border-gray-800 last:border-b-0">
                  <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-0">
                    <span className="text-base-content/70 text-xs">
                      [{formatTimestamp(log.timestamp)}]
                    </span>
                    <span className={`text-xs sm:ml-2 ${
                      log.success ? 'text-green-400' : 'text-red-400'
                    }`}>
                      {log.success ? '[OK]' : '[FAIL]'}
                    </span>
                  </div>
                  <span className="text-base-content/70 text-xs md:text-sm block mt-1">
                    {log.message}
                  </span>
                </div>
              ))
            )}
            
          </div>
        </div>

        <div className="flex justify-between items-center mt-2 text-xs text-base-content/60">
          <div>
            Total: {logs.length} entries
          </div>
        </div>
      </div>
    </div>
  )
}

export default MonitoringLogCard