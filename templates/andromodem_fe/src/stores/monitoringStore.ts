import { create } from "zustand"

import type { MonitoringLog } from "@/types/monitoring"
import type { MonitoringStore } from "@/types/monitoringStore"

const MAX_LOGS = 100

export const useMonitoringStore = create<MonitoringStore>((set) => {
    return {
        logs: [],
        addLog: (log: MonitoringLog) => {
            set((state) => {
                const updatedLogs = [...state.logs, log]
                
                const finalLogs = updatedLogs.length > MAX_LOGS 
                    ? updatedLogs.slice(-MAX_LOGS) 
                    : updatedLogs
                
                return {
                    logs: finalLogs
                }
            })
        },
        clearLogs: () => {
            set({ logs: [] })
        }
    }
})