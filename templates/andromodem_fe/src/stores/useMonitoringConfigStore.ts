import { create } from 'zustand'

import type { MonitoringConfigStore } from '@/types/monitoringStore'


export const useMonitoringConfigStore = create<MonitoringConfigStore>((set) => ({
  monitoringConfig: null,
  setMonitoringConfig: (config) => set({ monitoringConfig: config }),
  clearMonitoringConfig: () => set({ monitoringConfig: null })
}))